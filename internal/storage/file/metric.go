package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/konstantin-kukharev/metrics/domain/entity"
	"github.com/konstantin-kukharev/metrics/internal"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"go.uber.org/zap"
)

const (
	stateInit    = "init"
	stateRunning = "running"
	stateStop    = "stop"
)

type key struct {
	t entity.MType
	n string
}

type addRequest struct {
	request  []*entity.Metric
	response chan<- []*entity.Metric
}

type MetricStorage struct {
	state string
	log   *logger.Logger
	store map[key]*entity.Metric
	mx    *sync.RWMutex

	sourcePath    string // path to file storage
	source        io.Writer
	storeInterval time.Duration // time between storage updates
	restore       bool          // if true - restore metrics from file storage

	add chan addRequest
}

func (ms *MetricStorage) Set(ctx context.Context, es ...*entity.Metric) ([]*entity.Metric, error) {
	if ms.state != stateRunning {
		return nil, fmt.Errorf("storage is stopped for new connections")
	}

	resp := make(chan []*entity.Metric)
	req := new(addRequest)
	req.request = es
	req.response = resp

	ms.add <- *req

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case resp, ok := <-resp:
			if ok {
				return resp, nil
			}
			return nil, fmt.Errorf("error while adding")
		}
	}
}

func (ms *MetricStorage) Get(ctx context.Context, ems ...*entity.Metric) ([]*entity.Metric, bool) {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	result := make([]*entity.Metric, 0, len(ems))

	for _, m := range ems {
		k := key{t: m.MType, n: m.ID}
		if v, ok := ms.store[k]; ok {
			result = append(result, v)
		} else {
			return result, false
		}
	}

	return result, true
}

func (ms *MetricStorage) List(ctx context.Context) []*entity.Metric {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	list := make([]*entity.Metric, 0, len(ms.store))
	for _, val := range ms.store {
		list = append(list, val)
	}

	return list
}

func (ms *MetricStorage) update(_ context.Context, req addRequest) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	keysForSelect := make(map[struct {
		k entity.MType
		t string
	}]*entity.Metric)

	for _, m := range req.request {
		k := key{t: m.MType, n: m.ID}
		res, ok := ms.store[k]
		if ok {
			m.Aggregate(res)
		}
		ms.store[k] = m
		keysForSelect[struct {
			k entity.MType
			t string
		}{m.MType, m.ID}] = m
	}

	result := make([]*entity.Metric, 0, len(keysForSelect))
	for _, m := range keysForSelect {
		result = append(result, m)
	}

	req.response <- result
	close(req.response)
}

func (ms *MetricStorage) report(ctx context.Context, req addRequest) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	defer close(req.response)

	for _, m := range req.request {
		b, err := json.Marshal(m)
		if err != nil {
			ms.log.ErrorCtx(ctx, "marshal error", zap.Any("error", err))
			return
		}

		if _, err = ms.source.Write(b); err != nil {
			ms.log.ErrorCtx(ctx, "file write error", zap.Any("error", err))
			return
		}
		if _, err = ms.source.Write([]byte("\n")); err != nil {
			ms.log.ErrorCtx(ctx, "file write error", zap.Any("error", err))
			return
		}
	}

	req.response <- req.request
}

func (ms *MetricStorage) initialize(ctx context.Context) error {
	if ms.restore {
		ms.mx.Lock()
		ms.log.InfoCtx(ctx, "TRY TO RESTORE")
		restoreFile, err := os.OpenFile(ms.sourcePath, os.O_RDONLY|os.O_CREATE, internal.DefaultFileStoragePermission)
		if err != nil {
			ms.log.ErrorCtx(ctx, "RESTORE ERROR", zap.Any("error", err))
			return err
		}
		sc := bufio.NewScanner(restoreFile)
		for sc.Scan() {
			data := sc.Bytes()
			z := new(entity.Metric)
			if err := json.Unmarshal(data, z); err != nil {
				continue
			}
			k := key{t: z.MType, n: z.ID}
			ms.store[k] = z
		}
		restoreFile.Close()
		ms.mx.Unlock()
	}

	return nil
}

func (ms *MetricStorage) reporter(ctx context.Context) {
	for {
		select {
		case <-time.After(ms.storeInterval):
			c := context.WithoutCancel(ctx)
			res := make(chan []*entity.Metric)
			ms.report(c, addRequest{
				request:  ms.List(c),
				response: res,
			})
			<-res
		case <-ctx.Done():
			return
		}
	}
}

func (ms *MetricStorage) Run(ctx context.Context) error {
	if err := ms.initialize(ctx); err != nil {
		return err
	}

	file, err := os.OpenFile(ms.sourcePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, internal.DefaultFileStoragePermission)
	if err != nil {
		return err
	}
	defer file.Close()

	ms.source = file
	ms.state = stateRunning

	ms.log.InfoCtx(ctx, "file storage is running",
		zap.String("path", ms.sourcePath),
		zap.Duration("interval", ms.storeInterval),
		zap.Bool("restore", ms.restore))

	if ms.storeInterval != 0 {
		go ms.reporter(ctx)
	}

	for {
		select {
		case req := <-ms.add:
			c := context.WithoutCancel(ctx)
			if ms.storeInterval != 0 {
				ms.update(c, req)
			} else {
				res := make(chan []*entity.Metric, 1)
				ms.update(c, addRequest{
					request:  req.request,
					response: res,
				})
				resUpd, ok := <-res
				if !ok {
					close(req.response)
					break
				}
				ms.report(c, addRequest{
					request:  resUpd,
					response: req.response,
				})
			}
		case <-ctx.Done():
			ms.state = stateStop
			c := context.WithoutCancel(ctx)
			res := make(chan []*entity.Metric)
			ms.report(c, addRequest{
				request:  ms.List(c),
				response: res,
			})
			<-res

			return nil
		}
	}
}

func NewMetric(l *logger.Logger, restore bool, sourcePath string, storeInterval time.Duration) *MetricStorage {
	return &MetricStorage{
		log:           l,
		store:         make(map[key]*entity.Metric),
		mx:            &sync.RWMutex{},
		add:           make(chan addRequest),
		state:         stateInit,
		restore:       restore,
		sourcePath:    sourcePath,
		storeInterval: storeInterval,
	}
}

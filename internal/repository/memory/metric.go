package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/konstantin-kukharev/metrics/domain/entity"
	"github.com/konstantin-kukharev/metrics/internal/logger"
)

const (
	stateInit    = "init"
	stateRunning = "running"
	stateStop    = "stop"
)

type key struct {
	t, n string
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
			return nil, fmt.Errorf("error while adding ")
		}
	}
}

func (ms *MetricStorage) Get(ctx context.Context, m *entity.Metric) (*entity.Metric, bool) {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	k := key{t: m.MType, n: m.ID}
	if v, ok := ms.store[k]; ok {
		return v, ok
	}

	return m, false
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

func (ms *MetricStorage) update(ctx context.Context, req addRequest) {
	ms.mx.Lock()

	for _, m := range req.request {
		res, ok := ms.Get(ctx, m)
		if ok {
			m.Aggregate(res)
		}
		k := key{t: m.MType, n: m.ID}
		ms.store[k] = m
	}

	ms.mx.Unlock()
	req.response <- req.request
}

func (ms *MetricStorage) Run(ctx context.Context) error {
	ms.state = stateRunning
	for {
		select {
		case req := <-ms.add:
			c, cncl := context.WithTimeout(ctx, 1*time.Second)
			ms.update(c, req)
			cncl()
		case <-ctx.Done():
			ms.state = stateStop
			for req := range ms.add {
				c, cncl := context.WithTimeout(ctx, 1*time.Second)
				ms.update(c, req)
				cncl()
			}

			return nil
		}
	}
}

func NewMetric(l *logger.Logger) *MetricStorage {
	ms := new(MetricStorage)
	ms.log = l
	ms.store = map[key]*entity.Metric{}
	ms.mx = &sync.RWMutex{}
	ms.add = make(chan addRequest)
	ms.state = stateInit

	return ms
}

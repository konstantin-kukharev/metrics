package application

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/konstantin-kukharev/metrics/domain/entity"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"go.uber.org/zap"
)

type storage interface {
	List(context.Context) []*entity.Metric
}

type Reporter struct {
	cli *http.Client
	url string
	s   storage
	i   time.Duration
	log *logger.Logger
}

var retryIntervals = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

func NewReporter(l *logger.Logger, f *http.Client, s storage, url string, i time.Duration) *Reporter {
	return &Reporter{
		cli: f,
		url: url,
		s:   s,
		i:   i,
		log: l,
	}
}

func (r *Reporter) report(ctx context.Context) {
	b, err := json.Marshal(r.s.List(ctx))
	if err != nil {
		r.log.WarnCtx(ctx, "error while marshaling metric",
			zap.String("message", err.Error()),
		)
		return
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, r.url, bytes.NewBuffer(b))
	if err != nil {
		r.log.WarnCtx(ctx, "error while creating request",
			zap.String("message", err.Error()),
		)
		return
	}
	request.Header.Add("Content-Type", "application/json")
	r.requestRetry(ctx, request, retryIntervals...)
}

func (r *Reporter) requestRetry(ctx context.Context, req *http.Request, wait ...time.Duration) {
	intervals := make(chan struct{})

	go func(ctx context.Context, c chan<- struct{}, w []time.Duration) {
		for _, in := range w {
			select {
			case <-time.After(in):
				c <- struct{}{}
			case <-ctx.Done():
				close(c)
				return
			}
		}
		close(c)
	}(ctx, intervals, wait)

	for {
		select {
		case _, ok := <-intervals:
			if !ok {
				r.log.WarnCtx(ctx, "error while creating request",
					zap.String("message", "all report retry attempts are exhausted"),
				)

				return
			}
			resp, err := r.cli.Do(req)
			if err != nil {
				r.log.WarnCtx(ctx, "error while sending report",
					zap.String("message", err.Error()),
				)

				continue
			}
			resp.Body.Close()

			return
		case <-ctx.Done():
			r.log.WarnCtx(ctx, "error while sending report",
				zap.String("message", ctx.Err().Error()),
			)
			return
		}
	}
}

func (r *Reporter) Run(ctx context.Context) error {
	r.log.InfoCtx(ctx, "agent reporter is running")
	for {
		select {
		case <-time.After(r.i):
			r.log.InfoCtx(ctx, "new report")
			c := context.WithoutCancel(ctx)
			r.report(c)
		case <-ctx.Done():

			return nil
		}
	}
}

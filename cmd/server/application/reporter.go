package application

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/konstantin-kukharev/metrics/domain/entity"
	"github.com/konstantin-kukharev/metrics/domain/event"
)

type Reporter struct {
	f    io.Writer
	s    <-chan event.MetricAdd
	i    time.Duration
	data []*entity.Metric
}

func NewReporter(f io.Writer, s <-chan event.MetricAdd, i time.Duration) *Reporter {
	return &Reporter{
		f:    f,
		s:    s,
		i:    i,
		data: []*entity.Metric{},
	}
}

func (r *Reporter) report() {
	for _, m := range r.data {
		if b, err := json.Marshal(m); err == nil {
			_, _ = r.f.Write(b)
			_, _ = r.f.Write([]byte("\n"))
		}
	}
}

func (r *Reporter) Run(ctx context.Context) error {
	for {
		select {
		case em, ok := <-r.s:
			if !ok {
				return errors.New("add metric listener chan closed")
			}
			r.data = append(r.data, em.Metric)
		case <-time.After(r.i):
			r.report()
		case <-ctx.Done():
			r.report()

			return nil
		}
	}
}

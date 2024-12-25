package application

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/konstantin-kukharev/metrics/domain/event"
)

type StreamReporter struct {
	f io.Writer
	s <-chan event.MetricAdd
}

func NewStreamReporter(f io.Writer, s <-chan event.MetricAdd) *Reporter {
	return &Reporter{
		f: f,
		s: s,
	}
}

func (r *StreamReporter) Run(ctx context.Context) error {
	for {
		select {
		case em, ok := <-r.s:
			if !ok {
				return errors.New("add metric listener chan closed")
			}
			if b, err := json.Marshal(em.Metric); err == nil {
				_, _ = r.f.Write(b)
				_, _ = r.f.Write([]byte("\n"))
			}
		case <-ctx.Done():
			return nil
		}
	}
}

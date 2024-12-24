package application

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type Source interface {
	List() []*entity.Metric
}

type Reporter struct {
	f *os.File
	s Source
	i time.Duration
}

func NewReporter(f *os.File, s Source, i time.Duration) *Reporter {
	return &Reporter{
		f: f,
		s: s,
		i: i,
	}
}

func (r *Reporter) report() {
	for _, m := range r.s.List() {
		if b, err := json.Marshal(m); err == nil {
			_, _ = r.f.Write(b)
			_, _ = r.f.WriteString("\n")
		}
	}
}

func (r *Reporter) Run(ctx context.Context) error {
	for {
		select {
		case <-time.After(r.i):
			r.report()
		case <-ctx.Done():
			r.report()

			return nil
		}
	}
}

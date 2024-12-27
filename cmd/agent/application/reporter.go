package application

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type getter interface {
	Do() []*entity.Metric
}

type Reporter struct {
	cli *http.Client
	url string
	s   getter
	i   time.Duration
}

func NewReporter(f *http.Client, s getter, url string, i time.Duration) *Reporter {
	return &Reporter{
		cli: f,
		url: url,
		s:   s,
		i:   i,
	}
}

func (r *Reporter) report() {
	for _, m := range r.s.Do() {
		b, err := json.Marshal(m)
		if err != nil {
			continue
		}

		request, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, r.url, bytes.NewBuffer(b))
		if err != nil {
			continue
		}
		request.Header.Add("Content-Type", "application/json")
		resp, err := r.cli.Do(request)
		if err != nil {
			continue
		}
		resp.Body.Close()
	}
}

func (r *Reporter) Run(ctx context.Context) error {
	for {
		select {
		case <-time.After(r.i):
			r.report()
		case <-ctx.Done():

			return nil
		}
	}
}

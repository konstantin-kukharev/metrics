package application

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type storage interface {
	List(context.Context) []*entity.Metric
}

type Reporter struct {
	cli *http.Client
	url string
	s   storage
	i   time.Duration
}

func NewReporter(f *http.Client, s storage, url string, i time.Duration) *Reporter {
	return &Reporter{
		cli: f,
		url: url,
		s:   s,
		i:   i,
	}
}

func (r *Reporter) report(ctx context.Context) {
	for _, m := range r.s.List(ctx) {
		b, err := json.Marshal(m)
		if err != nil {
			continue
		}

		request, err := http.NewRequestWithContext(ctx, http.MethodPost, r.url, bytes.NewBuffer(b))
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
			c, cncl := context.WithDeadline(ctx, time.Now().Add(1*time.Second))
			r.report(c)
			cncl()
		case <-ctx.Done():

			return nil
		}
	}
}

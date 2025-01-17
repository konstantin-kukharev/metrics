package application

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/konstantin-kukharev/metrics/internal/logger"
	"github.com/konstantin-kukharev/metrics/internal/roundtripper"
	"go.uber.org/zap"
)

type Reporter struct {
	cli *http.Client
	url string
	s   storage
	i   time.Duration
	log *logger.Logger
}

func NewReporter(l *logger.Logger, s storage, url string, i time.Duration) *Reporter {
	var rt http.RoundTripper
	rt = http.DefaultTransport
	rt = roundtripper.NewRetry(rt, roundtripper.DefaultRetryDurations...)
	rt = roundtripper.NewCompress(rt)
	rt = roundtripper.NewLogging(rt, l)
	cli := &http.Client{
		Transport: rt,
		Timeout:   10 * time.Second,
	}

	return &Reporter{
		cli: cli,
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
	resp, err := r.cli.Do(request)
	if err != nil {
		r.log.WarnCtx(ctx, "error while sending report",
			zap.String("message", err.Error()),
		)

		return
	}
	resp.Body.Close()
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

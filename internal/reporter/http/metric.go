package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type Reporter struct {
	cli *http.Client
	url string
}

func (r *Reporter) Do(m *entity.Metric) error {
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, r.url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	resp, err := r.cli.Do(request)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func NewReporter(cli *http.Client, url string) *Reporter {
	return &Reporter{
		cli: cli,
		url: url,
	}
}

package report

import (
	"errors"
	"net/http"

	"github.com/konstantin-kukharev/metrics/internal/metric"
)

type AgentReporter interface {
	Report(serverAddress string, data []metric.Value) error
}
type Rest struct {
	cli *http.Client
}

func (r *Rest) Report(server string, d []metric.Value) error {
	var errs error
	for _, v := range d {
		val, err := v.GetValue()
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		url := "http://" + server + "/update/" + v.Type() + "/" + v.Name() + "/" + val
		res, err := r.cli.Post(url, "text/plain", http.NoBody)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		defer res.Body.Close()
	}

	return errs
}

func NewRest(cli *http.Client) *Rest {
	r := new(Rest)
	r.cli = cli

	return r
}

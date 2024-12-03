package report

import (
	"errors"
	"net/http"

	"github.com/konstantin-kukharev/metrics/pkg/metric"
)

type rest struct {
	cli *http.Client
}

func (r *rest) Report(server string, d []metric.Value) error {
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

func NewRest(cli *http.Client) *rest {
	return &rest{cli: cli}
}

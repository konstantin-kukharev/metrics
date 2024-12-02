package report

import (
	"fmt"
	"net/http"

	"github.com/konstantin-kukharev/metrics/internal"
)

type rest struct {
	cli *http.Client
}

func (r *rest) Report(server string, d []internal.MetricValue) {
	for _, v := range d {
		url := "http://" + server + "/update/" + v.Type() + "/" + v.Name() + "/" + v.Value()
		res, err := r.cli.Post(url, "text/plain", http.NoBody)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer res.Body.Close()
		fmt.Println(res.StatusCode)
	}
}

func NewRest(cli *http.Client) *rest {
	return &rest{cli: cli}
}

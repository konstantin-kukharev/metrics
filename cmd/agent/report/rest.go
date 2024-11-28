package report

import (
	"fmt"
	"net/http"

	"github.com/konstantin-kukharev/metrics/internal"
)

type rest struct {
	cli    *http.Client
	server string
}

func (r *rest) Report(d []internal.MetricValue) {
	for _, v := range d {
		url := "http://" + r.server + "/update/" + v.Type() + "/" + v.Name() + "/" + v.Value()
		res, err := r.cli.Post(url, "text/plain", http.NoBody)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer res.Body.Close()
		fmt.Println(res.StatusCode) //14339
	}
}

func NewRest(cli *http.Client, server string) *rest {
	return &rest{cli: cli, server: server}
}

package handler

import (
	"fmt"
	"net/http"

	"github.com/konstantin-kukharev/metrics/cmd/server/service"
)

type metricIndex struct {
	service service.Metric
}

func (s *metricIndex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	body := "<ul>"
	for _, d := range s.service.List() {
		val, err := d.GetValue()
		if err != nil {
			body += fmt.Sprintf("<li>%s\t%s\t%s</li>", d.Type(), d.Name(), err.Error())
			continue
		}
		body += fmt.Sprintf("<li>%s\t%s\t%s</li>", d.Type(), d.Name(), val)
	}

	body += "</ul>"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func NewIndexMetric(srv service.Metric) *metricIndex {
	serv := &metricIndex{service: srv}
	return serv
}

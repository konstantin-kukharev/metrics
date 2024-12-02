package handler

import (
	"fmt"
	"net/http"

	"github.com/konstantin-kukharev/metrics/internal"
)

type metricIndex struct {
	service internal.MetricService
}

func (s *metricIndex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	var body string = "<ul>"
	for _, d := range s.service.List() {
		body += fmt.Sprintf("<li>%s\t%s\t%s</li>", d.Type(), d.Name(), d.Value())
	}

	body += "</ul>"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func NewIndexMetric(srv internal.MetricService) *metricIndex {
	serv := &metricIndex{service: srv}
	return serv
}

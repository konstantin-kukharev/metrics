package handler

import (
	"net/http"

	"github.com/konstantin-kukharev/metrics/cmd/server/service"
)

type MetricGet struct {
	service service.Metric
}

func (s *MetricGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	t := r.PathValue("type")
	n := r.PathValue("name")

	w.Header().Add("Content-Type", "text/plain")

	if n == "" || t == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	if v, ok := s.service.Get(t, n); ok {
		w.Write(v)
		w.WriteHeader(http.StatusOK)
	}

	w.WriteHeader(http.StatusNotFound)
}

func NewGetMetric(srv service.Metric) *MetricGet {
	serv := &MetricGet{service: srv}
	return serv
}

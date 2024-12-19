package metric

import (
	"html/template"
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type MetricListReader interface {
	// List get all metric list
	Do() []*entity.Metric
}

type metricIndex struct {
	service MetricListReader
}

func (s *metricIndex) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	Data := s.service.Do()
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, Data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NewIndexMetric(srv MetricListReader) *metricIndex {
	serv := &metricIndex{service: srv}
	return serv
}

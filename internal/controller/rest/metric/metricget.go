package metric

import (
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain/entity"
	"github.com/konstantin-kukharev/metrics/internal"
)

type MetricReader interface {
	//Get get metric
	Do(m *entity.Metric) (*entity.Metric, bool)
}

type MetricGet struct {
	service MetricReader
}

func (s *MetricGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	t := r.PathValue("type")
	n := r.PathValue("name")

	w.Header().Add("Content-Type", "text/plain")

	if n == "" || t == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	data, err := entity.NewMetric(n, t, "")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if v, ok := s.service.Do(data); ok {
		res := internal.GetMetricValue(v)
		w.Write([]byte(res))
		w.WriteHeader(http.StatusOK)
	}

	w.WriteHeader(http.StatusNotFound)
}

func NewGetMetric(srv MetricReader) *MetricGet {
	serv := &MetricGet{service: srv}
	return serv
}

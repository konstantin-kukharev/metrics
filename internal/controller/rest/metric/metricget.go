package metric

import (
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type MetricReader interface {
	// Get get metric
	Do(m *entity.Metric) (*entity.Metric, bool)
}

type MetricGet struct {
	service MetricReader
}

func (s *MetricGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

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
		res := v.GetValue()
		_, err := w.Write([]byte(res))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}

	w.WriteHeader(http.StatusNotFound)
}

func NewGetMetric(srv MetricReader) *MetricGet {
	serv := &MetricGet{service: srv}
	return serv
}

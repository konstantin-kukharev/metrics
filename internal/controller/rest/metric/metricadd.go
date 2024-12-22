package metric

import (
	"errors"
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type MetricWriter interface {
	// Do set metric value
	Do(...*entity.Metric) error
}

type MetricAdd struct {
	service MetricWriter
}

func (s *MetricAdd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	t := r.PathValue("type")
	n := r.PathValue("name")
	v := r.PathValue("val")

	if v == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Add("Content-Type", "text/plain")

	data, err := entity.NewMetric(n, t, v)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWrongMetricName):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, domain.ErrWrongMetricType):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	if err := s.service.Do(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func NewAddMetric(srv MetricWriter) *MetricAdd {
	serv := &MetricAdd{service: srv}
	return serv
}

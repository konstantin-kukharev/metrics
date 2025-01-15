package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type metricWriter interface {
	Set(context.Context, ...*entity.Metric) ([]*entity.Metric, error)
}

type MetricAdd struct {
	service metricWriter
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
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err = data.Validate(); err != nil {
		switch {
		case errors.Is(err, domain.ErrWrongMetricName):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	if _, err := s.service.Set(r.Context(), data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewAddMetric(srv metricWriter) *MetricAdd {
	serv := &MetricAdd{service: srv}
	return serv
}

package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type metricReader interface {
	Get(context.Context, ...*entity.Metric) ([]*entity.Metric, bool)
}

type MetricGet struct {
	service metricReader
}

func (s *MetricGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	t := r.PathValue("type")
	n := r.PathValue("name")

	w.Header().Add("Content-Type", "text/plain")

	data, err := entity.NewMetric(n, t, "")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err = data.Validate(); err != nil && !errors.Is(err, domain.ErrEmptyMetricValue) {
		switch {
		case errors.Is(err, domain.ErrWrongMetricName):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	v, ok := s.service.Get(r.Context(), data)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(v) != 1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := v[0].GetValue()
	_, err = w.Write([]byte(res))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewGetMetric(srv metricReader) *MetricGet {
	serv := &MetricGet{service: srv}
	return serv
}

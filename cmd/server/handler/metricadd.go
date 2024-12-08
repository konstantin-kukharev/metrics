package handler

import (
	"errors"
	"net/http"

	"github.com/konstantin-kukharev/metrics/cmd/server/service"
	"github.com/konstantin-kukharev/metrics/domain"
)

type MetricAdd struct {
	service service.Metric
}

// - Прием метрик по протоколу НТТР методом POST
//
// - При успешном приёме возвращать http.StatusOk
//
// - При попытке передать запрос без имени метрики возвращать http.StatusNotFound
//
// - При попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest
//
// - Редиректы не поддерживаются.
func (s *MetricAdd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	t := r.PathValue("type")
	n := r.PathValue("name")
	v := r.PathValue("val")

	w.Header().Add("Content-Type", "text/plain")

	if n == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	if err := s.service.Set(t, n, v); err != nil {
		if errors.Is(err, domain.ErrInvalidData) || errors.Is(err, domain.ErrWrongMetricType) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func NewAddMetric(srv service.Metric) *MetricAdd {
	serv := &MetricAdd{service: srv}
	return serv
}

package main

import (
	"errors"
	"net/http"

	"github.com/konstantin-kukharev/metrics/cmd/server/internal"
)

type server struct {
	service internal.MetricService
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
func (s *server) MetricUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	t := r.PathValue("type")
	n := r.PathValue("name")
	v := r.PathValue("val")

	w.Header().Add("Content-Type", "text/plain")

	if n == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := s.service.Set(t, n, v); err != nil {
		if errors.Is(err, internal.ErrInvalidData) || errors.Is(err, internal.ErrWrongMetricType) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewServer(srv internal.MetricService) *server {
	serv := &server{service: srv}
	return serv
}

package handler

import (
	"context"
	"html/template"
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type metricListReader interface {
	List(context.Context) []*entity.Metric
}

type metricIndex struct {
	service metricListReader
}

func (s *metricIndex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	Data := s.service.List(r.Context())
	tpath := "site/index.html"

	tmpl, err := template.ParseFiles(tpath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, Data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NewIndexMetric(srv metricListReader) *metricIndex {
	serv := &metricIndex{service: srv}
	return serv
}

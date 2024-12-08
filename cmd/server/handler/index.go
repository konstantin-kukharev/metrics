package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/konstantin-kukharev/metrics/cmd/server/service"
)

type metricIndex struct {
	service service.Metric
}

func (s *metricIndex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("Content-Type", "text/html")
	Data := s.service.List()
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		fmt.Println(err)
	} else {
		tmpl.Execute(w, Data)
	}
}

func NewIndexMetric(srv service.Metric) *metricIndex {
	serv := &metricIndex{service: srv}
	return serv
}

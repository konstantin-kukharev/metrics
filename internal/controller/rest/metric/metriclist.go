package metric

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type MetricListReader interface {
	//List get all metric list
	List() []*entity.Metric
}

type metricIndex struct {
	service MetricListReader
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

func NewIndexMetric(srv MetricListReader) *metricIndex {
	serv := &metricIndex{service: srv}
	return serv
}

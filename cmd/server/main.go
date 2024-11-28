package main

import (
	"net/http"

	"log"

	"github.com/konstantin-kukharev/metrics/cmd/server/handler"
	"github.com/konstantin-kukharev/metrics/cmd/server/metric"
	"github.com/konstantin-kukharev/metrics/cmd/server/service"
	"github.com/konstantin-kukharev/metrics/cmd/server/storage"
	"github.com/konstantin-kukharev/metrics/internal"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

/*
run
*/
func run() error {
	store := storage.NewMemStorage()

	serv := service.NewMetric(store, metric.Gauge(), metric.Counter())
	srv := handler.NewMetric(serv)

	mux := http.NewServeMux()
	LinkMetric := "/update/{type}/{name}/{val}"
	mux.HandleFunc(LinkMetric, srv.MetricUpdate)

	return http.ListenAndServe(internal.DefaultServerAddr, mux)
}

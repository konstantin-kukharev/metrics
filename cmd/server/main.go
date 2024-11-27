package main

import (
	"net/http"

	"log"

	"github.com/konstantin-kukharev/metrics/cmd/server/metric"
	"github.com/konstantin-kukharev/metrics/cmd/server/storage"
)

var serverAddr = ":8080"

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

	serv := NewMetricService(store, metric.Gauge(), metric.Counter())
	srv := NewServer(serv)

	mux := http.NewServeMux()
	LinkMetric := "/update/{type}/{name}/{val}"
	mux.HandleFunc(LinkMetric, srv.MetricUpdate)

	return http.ListenAndServe(serverAddr, mux)
}

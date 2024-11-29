package main

import (
	"net/http"

	"log"

	"github.com/konstantin-kukharev/metrics/cmd/server/handler"
	"github.com/konstantin-kukharev/metrics/cmd/server/metric"
	"github.com/konstantin-kukharev/metrics/cmd/server/service"
	"github.com/konstantin-kukharev/metrics/cmd/server/storage"
	"github.com/konstantin-kukharev/metrics/internal"

	"github.com/go-chi/chi/v5"
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

	r := chi.NewRouter()
	r.Method("POST", "/update/{type}/{name}/{val}", handler.NewAddMetric(serv))
	r.Method("GET", "/value/{type}/{name}", handler.NewGetMetric(serv))
	mux := http.NewServeMux()

	return http.ListenAndServe(internal.DefaultServerAddr, mux)
}

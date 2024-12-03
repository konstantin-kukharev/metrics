package main

import (
	"net/http"

	"log"

	"github.com/konstantin-kukharev/metrics/cmd/server/handler"
	"github.com/konstantin-kukharev/metrics/cmd/server/service"
	"github.com/konstantin-kukharev/metrics/cmd/server/settings"
	"github.com/konstantin-kukharev/metrics/cmd/server/storage"
	"github.com/konstantin-kukharev/metrics/pkg/metric"

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
	conf := settings.New()
	store := storage.NewMemStorage()
	serv := service.NewMetric(conf, store, &metric.Gauge{}, &metric.Counter{})

	r := chi.NewRouter()
	r.Method("POST", "/update/{type}/{name}/{val}", handler.NewAddMetric(serv))
	r.Method("GET", "/value/{type}/{name}", handler.NewGetMetric(serv))
	r.Method("GET", "/", handler.NewIndexMetric(serv))

	err := http.ListenAndServe(conf.GetAddress(), r)
	// fmt.Printf(
	// 	"runninig server on \"%s\"\r\n",
	// 	conf.GetAddress(),
	// )

	return err
}

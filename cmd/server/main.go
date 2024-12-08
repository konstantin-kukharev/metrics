package main

import (
	"fmt"
	"net/http"

	"github.com/konstantin-kukharev/metrics/cmd/server/handler"
	"github.com/konstantin-kukharev/metrics/cmd/server/service"
	"github.com/konstantin-kukharev/metrics/cmd/server/settings"
	"github.com/konstantin-kukharev/metrics/cmd/server/storage"
	"github.com/konstantin-kukharev/metrics/internal/metric"

	"github.com/go-chi/chi/v5"
)

func main() {
	conf := settings.NewConfig().WithFlag().WithEnv().WithDebug()

	if err := run(conf); err != nil {
		conf.Log().Error("error occured", "error", err)
	}
}

/*
run
*/
func run(app settings.Application) error {

	store := storage.NewMemStorage()
	serv := service.NewMetric(app.Log(), store, &metric.Gauge{}, &metric.Counter{})

	r := chi.NewRouter()
	r.Method("POST", "/update/{type}/{name}/{val}", handler.NewAddMetric(serv))
	r.Method("GET", "/value/{type}/{name}", handler.NewGetMetric(serv))
	r.Method("GET", "/", handler.NewIndexMetric(serv))

	fmt.Printf(
		"runninig server on \"%s\"\r\n",
		app.GetAddress(),
	)

	err := http.ListenAndServe(app.GetAddress(), r)

	return err
}

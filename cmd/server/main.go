package main

import (
	"fmt"
	"net/http"
	"time"

	ucase "github.com/konstantin-kukharev/metrics/cmd/domain/usecase/metric"
	handler "github.com/konstantin-kukharev/metrics/cmd/internal/controller/rest/metric"
	"github.com/konstantin-kukharev/metrics/cmd/internal/logger"
	"github.com/konstantin-kukharev/metrics/cmd/server/settings"
	"github.com/konstantin-kukharev/metrics/internal/repository/memory"

	"github.com/go-chi/chi/v5"
)

type ApplicationConfig interface {
	GetAddress() string
}

type Logger interface {
	Info(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Error(msg string, fields ...any)
}

func main() {
	conf := settings.NewConfig().WithFlag().WithEnv()
	log := logger.NewSlog()

	if err := run(conf, log); err != nil {
		log.Error("error occured", "error", err)
	}
}

/*
run
*/
func run(app ApplicationConfig, l Logger) error {

	store := memory.NewStorage(l)
	add := ucase.NewAddMetric(store)
	getVal := ucase.NewGetMetric(store)
	list := ucase.NewListMetric(store)

	r := chi.NewRouter()
	r.Method("POST", "/update/{type}/{name}/{val}", WithLogging(handler.NewAddMetric(add), l))
	r.Method("GET", "/value/{type}/{name}", WithLogging(handler.NewGetMetric(getVal), l))
	r.Method("GET", "/", WithLogging(handler.NewIndexMetric(list), l))

	fmt.Printf(
		"runninig server on \"%s\"\r\n",
		app.GetAddress(),
	)

	err := http.ListenAndServe(app.GetAddress(), r)

	return err
}

func WithLogging(h http.Handler, l Logger) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(w, r)

		duration := time.Since(start)
		l.Debug("new request",
			"uri", uri,
			"method", method,
			"duration", duration,
		)
	}

	return http.HandlerFunc(logFn)
}

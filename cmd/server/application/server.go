package application

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/konstantin-kukharev/metrics/internal/handler"
	"github.com/konstantin-kukharev/metrics/internal/middleware"
)

type ApplicationConfig interface {
	GetAddress() string
}

type Logger interface {
	Info(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
}

type Server struct {
	config ApplicationConfig
	log    Logger
	server *http.Server
}

func NewServer(
	w handler.MetricWriter,
	r handler.MetricReader,
	lr handler.MetricListReader,
	app ApplicationConfig,
	l Logger) *Server {
	router := chi.NewRouter()
	router.Method("POST", "/update/{type}/{name}/{val}", middleware.WithLogging(handler.NewAddMetric(w), l))
	router.Method("GET", "/value/{type}/{name}", middleware.WithLogging(handler.NewGetMetric(r), l))
	router.Method("GET", "/", middleware.WithCompressing(middleware.WithLogging(handler.NewIndexMetric(lr), l)))

	router.Method("POST", "/update/", middleware.WithCompressing(middleware.WithLogging(handler.NewAddMetricV2(w), l)))
	router.Method("POST", "/value/", middleware.WithCompressing(middleware.WithLogging(handler.NewMetricGetV2(r), l)))

	return &Server{
		config: app,
		log:    l,
		server: &http.Server{
			Handler:           router,
			Addr:              app.GetAddress(),
			ReadHeaderTimeout: 1 * time.Second,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func(c context.Context) {
		<-c.Done()

		err := s.server.Shutdown(c)
		if err != nil {
			s.log.Error("http server shutting down: " + err.Error())
		} else {
			s.log.Info("http server shutdown processed successfully")
		}
	}(ctx)

	return s.server.ListenAndServe()
}

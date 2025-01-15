package application

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/konstantin-kukharev/metrics/internal/handler"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"github.com/konstantin-kukharev/metrics/internal/middleware"
	"github.com/konstantin-kukharev/metrics/internal/storage"

	"go.uber.org/zap"
)

type ApplicationConfig interface {
	GetAddress() string
}

type Server struct {
	address string
	log     *logger.Logger
	server  *http.Server
}

func NewServer(
	l *logger.Logger,
	s storage.Metric,
	address string,
	databaseDNS string) *Server {
	router := chi.NewRouter()
	router.Method("POST", "/update/{type}/{name}/{val}", middleware.WithLogging(handler.NewAddMetric(s), l))
	router.Method("GET", "/value/{type}/{name}", middleware.WithLogging(handler.NewGetMetric(s), l))
	router.Method("GET", "/", middleware.WithCompressing(middleware.WithLogging(handler.NewIndexMetric(s), l)))

	router.Method("POST", "/update/", middleware.WithJSONContent(
		middleware.WithCompressing(
			middleware.WithLogging(
				handler.NewAddMetricV2(s), l))))
	router.Method("POST", "/updates/", middleware.WithJSONContent(
		middleware.WithCompressing(
			middleware.WithLogging(
				handler.NewAddMetricV3(s), l))))
	router.Method("POST", "/value/", middleware.WithJSONContent(
		middleware.WithCompressing(
			middleware.WithLogging(
				handler.NewMetricGetV2(s), l))))

	router.Method("GET", "/ping", middleware.WithLogging(handler.NewPing(databaseDNS, l), l))

	return &Server{
		address: address,
		log:     l,
		server: &http.Server{
			ErrorLog:          l.Std(),
			Handler:           router,
			Addr:              address,
			ReadHeaderTimeout: 1 * time.Second,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func(c context.Context) {
		<-c.Done()

		err := s.server.Shutdown(c)
		if err != nil {
			s.log.ErrorCtx(c, "http server shutting down",
				zap.String("error", err.Error()),
			)
		} else {
			s.log.InfoCtx(c, "http server shutdown processed successfully")
		}
	}(ctx)

	s.log.InfoCtx(ctx, "http server running",
		zap.String("address", s.address),
	)
	return s.server.ListenAndServe()
}

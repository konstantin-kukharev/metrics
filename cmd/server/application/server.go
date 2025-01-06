package application

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/konstantin-kukharev/metrics/domain/entity"
	"github.com/konstantin-kukharev/metrics/internal/handler"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"github.com/konstantin-kukharev/metrics/internal/middleware"

	"go.uber.org/zap"
)

type ApplicationConfig interface {
	GetAddress() string
	GetDatabaseDNS() string
}

type Server struct {
	config ApplicationConfig
	log    *logger.Logger
	server *http.Server
}

type repo interface {
	Set(context.Context, ...*entity.Metric) ([]*entity.Metric, error)
	Get(context.Context, *entity.Metric) (*entity.Metric, bool)
	List(context.Context) []*entity.Metric
}

func NewServer(
	l *logger.Logger,
	s repo,
	app ApplicationConfig) *Server {
	router := chi.NewRouter()
	router.Method("POST", "/update/{type}/{name}/{val}", middleware.WithLogging(handler.NewAddMetric(s), l))
	router.Method("GET", "/value/{type}/{name}", middleware.WithLogging(handler.NewGetMetric(s), l))
	router.Method("GET", "/", middleware.WithCompressing(middleware.WithLogging(handler.NewIndexMetric(s), l)))

	router.Method("POST", "/update/", middleware.WithJSONContent(
		middleware.WithCompressing(
			middleware.WithLogging(
				handler.NewAddMetricV2(s), l))))
	router.Method("POST", "/value/", middleware.WithJSONContent(
		middleware.WithCompressing(
			middleware.WithLogging(
				handler.NewMetricGetV2(s), l))))

	router.Method("GET", "/ping", middleware.WithLogging(handler.NewPing(app.GetDatabaseDNS(), l), l))

	return &Server{
		config: app,
		log:    l,
		server: &http.Server{
			ErrorLog:          l.Std(),
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
			s.log.ErrorCtx(c, "http server shutting down",
				zap.String("error", err.Error()),
			)
		} else {
			s.log.InfoCtx(c, "http server shutdown processed successfully")
		}
	}(ctx)

	return s.server.ListenAndServe()
}

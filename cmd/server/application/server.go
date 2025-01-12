package application

import (
	"context"
	"net/http"
	"time"

	"github.com/konstantin-kukharev/metrics/internal/logger"

	"go.uber.org/zap"
)

type ApplicationConfig interface {
	GetAddress() string
}

type Server struct {
	config ApplicationConfig
	log    *logger.Logger
	server *http.Server
}

func NewServer(
	l *logger.Logger,
	router http.Handler,
	app ApplicationConfig) *Server {
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

	s.log.InfoCtx(ctx, "http server running",
		zap.String("address", s.config.GetAddress()),
	)
	return s.server.ListenAndServe()
}

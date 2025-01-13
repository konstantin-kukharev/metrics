package main

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/konstantin-kukharev/metrics/cmd/server/application"
	"github.com/konstantin-kukharev/metrics/cmd/server/settings"
	"github.com/konstantin-kukharev/metrics/internal/graceful"
	"github.com/konstantin-kukharev/metrics/internal/handler"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"github.com/konstantin-kukharev/metrics/internal/middleware"
	"github.com/konstantin-kukharev/metrics/internal/storage"
	"github.com/konstantin-kukharev/metrics/internal/storage/file"
	"github.com/konstantin-kukharev/metrics/internal/storage/memory"
	"github.com/konstantin-kukharev/metrics/internal/storage/persistence"
	"go.uber.org/zap"
)

func main() {
	conf := settings.NewConfig()
	conf.WithFlag()
	conf.WithEnv()

	ctx := context.Background()
	l, err := logger.NewLogger(zap.InfoLevel)
	if err != nil {
		log.Panic(err)
	}
	ctx = l.WithContextFields(ctx,
		zap.Int("pid", os.Getpid()),
		zap.String("app", "server"))
	defer l.Sync()

	l.InfoCtx(ctx, "server running with options", zap.Any("config", conf))

	var s storage.Metric
	s = memory.NewMetric(l)
	if conf.GetDatabaseDNS() != "" {
		s = persistence.NewMetric(l, conf.GetDatabaseDNS())
	} else if conf.FileStoragePath != "" {
		s = file.NewMetric(l, conf)
	}

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

	router.Method("GET", "/ping", middleware.WithLogging(handler.NewPing(conf.GetDatabaseDNS(), l), l))
	server := application.NewServer(l, router, conf)

	gs := graceful.NewGracefulShutdown(ctx, 1*time.Second)
	gs.AddTask(s)
	gs.AddTask(server)
	err = gs.Wait(syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		l.ErrorCtx(ctx, "server finished", zap.Any("error", err))
	}
}

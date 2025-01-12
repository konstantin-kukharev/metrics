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
	"github.com/konstantin-kukharev/metrics/internal/repository"
	"github.com/konstantin-kukharev/metrics/internal/repository/file"
	"github.com/konstantin-kukharev/metrics/internal/repository/memory"
	"github.com/konstantin-kukharev/metrics/internal/repository/persistence"
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

	var storage repository.Metric
	storage = memory.NewMetric(l)
	if conf.GetDatabaseDNS() != "" {
		storage = persistence.NewMetric(l, conf.GetDatabaseDNS())
	} else if conf.FileStoragePath != "" {
		storage = file.NewMetric(l, conf)
	}

	router := chi.NewRouter()
	router.Method("POST", "/update/{type}/{name}/{val}", middleware.WithLogging(handler.NewAddMetric(storage), l))
	router.Method("GET", "/value/{type}/{name}", middleware.WithLogging(handler.NewGetMetric(storage), l))
	router.Method("GET", "/", middleware.WithCompressing(middleware.WithLogging(handler.NewIndexMetric(storage), l)))

	router.Method("POST", "/update/", middleware.WithJSONContent(
		middleware.WithCompressing(
			middleware.WithLogging(
				handler.NewAddMetricV2(storage), l))))
	router.Method("POST", "/updates/", middleware.WithJSONContent(
		middleware.WithCompressing(
			middleware.WithLogging(
				handler.NewAddMetricV3(storage), l))))
	router.Method("POST", "/value/", middleware.WithJSONContent(
		middleware.WithCompressing(
			middleware.WithLogging(
				handler.NewMetricGetV2(storage), l))))

	router.Method("GET", "/ping", middleware.WithLogging(handler.NewPing(conf.GetDatabaseDNS(), l), l))
	server := application.NewServer(l, router, conf)

	gs := graceful.NewGracefulShutdown(ctx, 1*time.Second)
	gs.AddTask(storage)
	gs.AddTask(server)
	err = gs.Wait(syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		l.ErrorCtx(ctx, "server finished", zap.Any("error", err))
	}
}

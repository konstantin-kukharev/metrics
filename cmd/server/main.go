package main

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/konstantin-kukharev/metrics/cmd/server/application"
	"github.com/konstantin-kukharev/metrics/cmd/server/settings"
	"github.com/konstantin-kukharev/metrics/internal/graceful"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"github.com/konstantin-kukharev/metrics/internal/repository"
	"github.com/konstantin-kukharev/metrics/internal/repository/memory"
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

	var storage repository.Metric
	switch {
	case conf.GetDatabaseDNS() != "":
	case conf.FileStoragePath != "" && conf.GetStoreInterval() != 0:
	case conf.FileStoragePath != "" && conf.GetStoreInterval() == 0:
	default:
		storage = memory.NewMetric(l)
	}

	server := application.NewServer(l, storage, conf)

	gs := graceful.NewGracefulShutdown(ctx, 1*time.Second)
	gs.AddTask(storage)
	gs.AddTask(server)
	err = gs.Wait(syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		l.FatalCtx(ctx, err.Error())
	}
}

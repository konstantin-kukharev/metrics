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
	if conf.DatabaseDNS != "" {
		s = persistence.NewMetric(l, conf.DatabaseDNS)
	} else if conf.FileStoragePath != "" {
		s = file.NewMetric(l, conf.Restore, conf.FileStoragePath, time.Duration(conf.StoreInterval*int(time.Second)))
	}

	server := application.NewServer(l, s, conf.Address, conf.DatabaseDNS)

	gs := graceful.NewGracefulShutdown(ctx, 1*time.Second)
	gs.AddTask(s)
	gs.AddTask(server)
	err = gs.Wait(syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		l.ErrorCtx(ctx, "server finished", zap.Any("error", err))
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/konstantin-kukharev/metrics/cmd/agent/application"
	"github.com/konstantin-kukharev/metrics/cmd/agent/settings"
	"go.uber.org/zap"

	"github.com/konstantin-kukharev/metrics/internal/graceful"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"github.com/konstantin-kukharev/metrics/internal/storage/memory"
)

func main() {
	conf := settings.New().WithFlag().WithEnv()

	ctx := context.Background()
	l, err := logger.NewLogger(zap.InfoLevel)
	if err != nil {
		log.Panic(err)
	}
	ctx = l.WithContextFields(ctx,
		zap.Int("pid", os.Getpid()),
		zap.String("app", "server"))
	defer l.Sync()

	l.InfoCtx(ctx, "agent running with options", zap.Any("config", conf))

	store := memory.NewMetric(l)

	reporter := application.NewReporter(l, store,
		fmt.Sprintf("http://%s/updates/", conf.Address), time.Duration(conf.ReportInterval*int(time.Second)))
	agent := application.NewAgent(store, time.Duration(conf.PoolInterval*int(time.Second)), l)

	gs := graceful.NewGracefulShutdown(ctx, 1*time.Second)
	gs.AddTask(store)
	gs.AddTask(agent)
	gs.AddTask(reporter)

	err = gs.Wait(syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		l.ErrorCtx(ctx, "agent service finished", zap.Any("error", err))
	}
}

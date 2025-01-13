package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/konstantin-kukharev/metrics/cmd/agent/application"
	"github.com/konstantin-kukharev/metrics/cmd/agent/settings"
	"go.uber.org/zap"

	"github.com/konstantin-kukharev/metrics/internal/graceful"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"github.com/konstantin-kukharev/metrics/internal/roundtripper"
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

	var rt http.RoundTripper
	rt = http.DefaultTransport
	rt = roundtripper.NewRetry(rt, roundtripper.DefaultRetryDurations...)
	rt = roundtripper.NewCompress(rt)
	rt = roundtripper.NewLogging(rt, l)
	cli := &http.Client{
		Transport: rt,
		Timeout:   10 * time.Second,
	}
	reporter := application.NewReporter(l, cli, store,
		fmt.Sprintf("http://%s/updates/", conf.GetServerAddress()), conf.GetReportInterval())
	agent := application.NewAgent(store, conf, l)

	gs := graceful.NewGracefulShutdown(ctx, 1*time.Second)
	gs.AddTask(store)
	gs.AddTask(agent)
	gs.AddTask(reporter)

	err = gs.Wait(syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		l.ErrorCtx(ctx, "agent service finished", zap.Any("error", err))
	}
}

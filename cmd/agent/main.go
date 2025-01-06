package main

import (
	"context"
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
	"github.com/konstantin-kukharev/metrics/internal/repository/memory"
	"github.com/konstantin-kukharev/metrics/internal/roundtripper"
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

	store := memory.NewMetric(l)

	var rt http.RoundTripper
	rt = http.DefaultTransport
	rt = roundtripper.NewLogging(rt, l)
	rt = roundtripper.NewCompress(rt)
	cli := &http.Client{
		Transport: rt,
	}
	reporter := application.NewReporter(cli, store, "http://"+conf.GetServerAddress()+"/update/", conf.GetReportInterval())
	agent := application.NewAgent(store, conf, l.Std())

	gs := graceful.NewGracefulShutdown(ctx, 1*time.Second)
	gs.AddTask(store)
	gs.AddTask(agent)
	gs.AddTask(reporter)

	err = gs.Wait(syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		l.FatalCtx(ctx, err.Error())
	}
}

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

	usecase "github.com/konstantin-kukharev/metrics/domain/usecase/metric"
)

func main() {
	conf := settings.New().WithFlag().WithEnv()

	ctx := context.Background()
	l, err := logger.NewZapLogger(zap.InfoLevel)
	if err != nil {
		log.Panic(err)
	}
	ctx = l.WithContextFields(ctx,
		zap.Int("pid", os.Getpid()),
		zap.String("app", "agent"))

	defer l.Sync()

	gs := graceful.NewGracefulShutdown(ctx, 1*time.Second)

	store := memory.NewStorage(l)
	add := usecase.NewAddMetric(store, nil)
	get := usecase.NewListMetric(store)

	var rt http.RoundTripper
	rt = http.DefaultTransport
	rt = roundtripper.NewLogging(rt, l)
	rt = roundtripper.NewCompress(rt)
	cli := &http.Client{
		Transport: rt,
	}
	r := application.NewReporter(cli, get, "http://"+conf.GetServerAddress()+"/update/", conf.GetReportInterval())
	gs.AddTask(r)

	agent := application.NewAgent(add, conf, l.Std())
	gs.AddTask(agent)

	err = gs.Wait(syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		l.FatalCtx(ctx, err.Error())
	}
}

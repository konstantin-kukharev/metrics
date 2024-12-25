package main

import (
	"context"
	"net/http"
	"syscall"
	"time"

	"github.com/konstantin-kukharev/metrics/cmd/agent/application"
	"github.com/konstantin-kukharev/metrics/cmd/agent/settings"

	"github.com/konstantin-kukharev/metrics/internal/graceful"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"github.com/konstantin-kukharev/metrics/internal/repository/memory"
	"github.com/konstantin-kukharev/metrics/internal/roundtripper"

	usecase "github.com/konstantin-kukharev/metrics/domain/usecase/metric"
)

func main() {
	conf := settings.New().WithFlag().WithEnv()
	log := logger.NewSlog()
	log.WithDebug()
	ctx := context.WithoutCancel(context.Background())
	gs := graceful.NewGracefulShutdown(ctx, 1*time.Second)

	store := memory.NewStorage(log)
	add := usecase.NewAddMetric(store, nil)
	get := usecase.NewListMetric(store)

	var rt http.RoundTripper
	rt = http.DefaultTransport
	rt = roundtripper.NewLogging(rt, log)
	rt = roundtripper.NewCompress(rt)
	cli := &http.Client{
		Transport: rt,
	}
	r := application.NewReporter(cli, get, "http://"+conf.GetServerAddress()+"/update/", conf.GetReportInterval())
	gs.AddTask(r)

	agent := application.NewAgent(add, conf, log)
	gs.AddTask(agent)

	err := gs.Wait(syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		log.Error("error occurred", "error", err)
	}
}

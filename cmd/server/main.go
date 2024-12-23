package main

import (
	"context"
	"syscall"
	"time"

	"github.com/konstantin-kukharev/metrics/cmd/server/application"
	"github.com/konstantin-kukharev/metrics/cmd/server/settings"
	"github.com/konstantin-kukharev/metrics/internal/graceful"
	"github.com/konstantin-kukharev/metrics/internal/logger"
)

type Logger interface {
	Info(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Error(msg string, fields ...any)
}

func main() {
	conf := settings.NewConfig().WithFlag().WithEnv()
	log := logger.NewSlog()
	log.WithDebug()
	ctx := context.WithoutCancel(context.Background())

	gs := graceful.NewGracefulShutdown(ctx, 1*time.Second)
	gs.AddTask(application.NewServer(conf, log))
	err := gs.Wait(syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		log.Error("error occurred", "error", err)
	}
}

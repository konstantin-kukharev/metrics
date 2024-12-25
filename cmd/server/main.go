package main

import (
	"bufio"
	"context"
	"os"
	"syscall"
	"time"

	"encoding/json"

	"github.com/konstantin-kukharev/metrics/cmd/server/application"
	"github.com/konstantin-kukharev/metrics/cmd/server/settings"
	"github.com/konstantin-kukharev/metrics/domain/entity"
	"github.com/konstantin-kukharev/metrics/domain/event"
	usecase "github.com/konstantin-kukharev/metrics/domain/usecase/metric"
	"github.com/konstantin-kukharev/metrics/internal"
	"github.com/konstantin-kukharev/metrics/internal/graceful"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"github.com/konstantin-kukharev/metrics/internal/repository/memory"
	"github.com/konstantin-kukharev/metrics/internal/transport"
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

	updater := make(chan event.MetricAdd)
	bus := transport.NewEventBus(updater)
	gs.AddTask(bus)

	store := memory.NewStorage(log)
	add := usecase.NewAddMetric(store, bus)
	getVal := usecase.NewGetMetric(store)
	list := usecase.NewListMetric(store)
	if conf.GetRestore() {
		file, err := os.OpenFile(conf.GetFileStoragePath(), os.O_RDONLY|os.O_CREATE, internal.DefaultFileStoragePermission)
		if err != nil {
			log.Error("open file", "error", err)
			return
		}
		sc := bufio.NewScanner(file)
		for sc.Scan() {
			data := sc.Bytes()
			z := new(entity.Metric)
			if err := json.Unmarshal(data, z); err != nil {
				continue
			}
			store.Set(z)
		}
		file.Close()
	}

	file, err := os.OpenFile(conf.GetFileStoragePath(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, internal.DefaultFileStoragePermission)
	if err != nil {
		log.Error("open file", "error", err)
		return
	}
	defer file.Close()

	if conf.GetStoreInterval() == 0 {
		events := make(chan event.MetricAdd)
		bus.AddListener(events)
		report := application.NewStreamReporter(file, events)
		gs.AddTask(report)
	}

	serverTask := application.NewServer(add, getVal, list, conf, log)
	gs.AddTask(serverTask)

	if conf.GetStoreInterval() > 0 {
		events := make(chan event.MetricAdd)
		bus.AddListener(events)
		report := application.NewReporter(file, events, conf.GetStoreInterval())
		gs.AddTask(report)
	}

	err = gs.Wait(syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		log.Error("error occurred", "error", err)
	}
}

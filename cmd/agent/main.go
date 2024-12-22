package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/konstantin-kukharev/metrics/cmd/agent/settings"
	"github.com/konstantin-kukharev/metrics/domain/entity"
	ucase "github.com/konstantin-kukharev/metrics/domain/usecase/metric"
	"github.com/konstantin-kukharev/metrics/internal"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"github.com/konstantin-kukharev/metrics/internal/repository/memory"
)

type ApplicationConfig interface {
	GetServerAddress() string
	GetReportInterval() time.Duration
	GetPoolInterval() time.Duration
}

type Logger interface {
	Info(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Error(msg string, fields ...any)
}

func main() {
	app := settings.New().WithFlag().WithEnv()
	log := logger.NewSlog()

	if err := run(app, log); err != nil {
		log.Error("error occurred", "error", err)
	}
}

func run(app *settings.Config, l Logger) error {
	store := memory.NewStorage(l)
	state := NewRuntimeMetric()
	add := ucase.NewAddMetric(store)
	list := ucase.NewListMetric(store)
	nextPool := time.Now()
	nextReport := time.Now()
	cli := &http.Client{}

	time.Sleep(internal.DefaultPoolInterval * time.Second)

	for {
		cTime := time.Now()
		if nextPool.Before(cTime) || nextPool.Equal(cTime) {
			var mem runtime.MemStats
			for _, stat := range state.List(&mem) {
				err := add.Do(stat)
				if err != nil {
					l.Error("error while updating runtime metrics",
						"msg", err.Error(),
					)

					return err
				}
			}
			nextPool = cTime.Add(app.GetPoolInterval())
		}
		if nextReport.Before(cTime) || nextReport.Equal(cTime) {
			err := report(cli, app.GetServerAddress(), list.Do()...)
			if err != nil {
				return err
			}
			nextReport = cTime.Add(app.GetReportInterval())
		}
	}
}

func report(cli *http.Client, server string, d ...*entity.Metric) error {
	for _, v := range d {
		body, err := json.Marshal(v)
		if err != nil {
			return err
		}
		url := "http://" + server + "/update"
		request, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		request.Header.Add("Content-Type", "application/json")
		res, err := cli.Do(request)
		if err != nil {
			return err
		}
		res.Body.Close()
	}

	return nil
}

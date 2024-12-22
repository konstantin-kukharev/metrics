package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"runtime"
	"time"

	"github.com/konstantin-kukharev/metrics/cmd/agent/settings"
	"github.com/konstantin-kukharev/metrics/domain/entity"
	ucase "github.com/konstantin-kukharev/metrics/domain/usecase/metric"
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

	//time.Sleep(app.GetPoolInterval() * time.Second)

	for {
		cTime := time.Now()
		if nextPool.Before(cTime) || nextPool.Equal(cTime) {
			l.Info("update pool",
				"time", cTime,
			)
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)

			err := add.Do(state.List(&mem)...)
			if err != nil {
				l.Info("error while updating runtime metrics",
					"msg", err.Error(),
				)

				return err
			}
			nextPool = cTime.Add(app.GetPoolInterval())
		}
		if nextReport.Before(cTime) || nextReport.Equal(cTime) {
			cli := &http.Client{}
			r := list.Do()
			err := report(cli, app.GetServerAddress(), r...)
			if err != nil {
				l.Error("error while reporting runtime metrics", err.Error())
			} else {
				l.Info("REPORT SUCCESS")
			}
			nextReport = cTime.Add(app.GetReportInterval())
		}
	}
}

func report(cli *http.Client, server string, d ...*entity.Metric) error {
	var errs error
	for _, v := range d {
		body, err := json.Marshal(v)
		if err != nil {
			errors.Join(errs, err)
			continue
		}
		url := "http://" + server + "/update/"
		request, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			errors.Join(errs, err)
			continue
		}
		request.Header.Add("Content-Type", "application/json")
		_, err = cli.Do(request)
		if err != nil {
			errors.Join(errs, err)
			continue
		}
	}

	return errs
}

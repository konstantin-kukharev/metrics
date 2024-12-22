package main

import (
	"net/http"
	"runtime"
	"time"

	"github.com/konstantin-kukharev/metrics/cmd/agent/settings"
	ucase "github.com/konstantin-kukharev/metrics/domain/usecase/metric"
	"github.com/konstantin-kukharev/metrics/internal"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	httpReport "github.com/konstantin-kukharev/metrics/internal/reporter/http"
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

	cli := &http.Client{}
	rp := httpReport.NewReporter(cli, "http://"+app.GetServerAddress()+"/update/")
	reporter := ucase.NewReportMetric(store, rp)

	nextPool := time.Now()
	nextReport := time.Now()

	time.Sleep(internal.DefaultPoolInterval * time.Second)

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
			err := reporter.Do()
			if err != nil {
				l.Error("error while reporting runtime metrics", err.Error())
			} else {
				l.Info("REPORT SUCCESS")
			}
			nextReport = cTime.Add(app.GetReportInterval())
		}
	}
}

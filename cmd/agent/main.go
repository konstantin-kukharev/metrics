package main

import (
	"net/http"
	"time"

	"github.com/konstantin-kukharev/metrics/cmd/agent/report"
	"github.com/konstantin-kukharev/metrics/cmd/agent/service"
	"github.com/konstantin-kukharev/metrics/cmd/agent/settings"
	"github.com/konstantin-kukharev/metrics/cmd/agent/state"
	"github.com/konstantin-kukharev/metrics/internal"
)

func main() {
	app := settings.New().WithFlag().WithEnv().WithDebug()
	c := &http.Client{}
	s := state.NewMemory()
	r := report.NewRest(c)
	srv := service.NewState(app, s, r)

	time.Sleep(internal.DefaultPoolInterval * time.Second)

	if err := srv.Run(); err != nil {
		app.Log().Error(
			"runninig agent report",
			"address", app.GetServerAddress(),
			"report interval", app.GetReportInterval(),
			"pool interval", app.GetPoolInterval(),
		)

		app.Log().Error("error occured", "message", err.Error())
	}
}

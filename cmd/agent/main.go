package main

import (
	"net/http"

	"github.com/konstantin-kukharev/metrics/cmd/agent/report"
	"github.com/konstantin-kukharev/metrics/cmd/agent/service"
	"github.com/konstantin-kukharev/metrics/cmd/agent/state"
	"github.com/konstantin-kukharev/metrics/internal"
)

func main() {
	c := &http.Client{}
	s := state.NewMemory()
	r := report.NewRest(c, internal.DefaultServerAddr)
	srv := service.NewState(s, r,
		internal.DefaultPoolInterval, internal.DefaultReportInterval)
	srv.Run()
}

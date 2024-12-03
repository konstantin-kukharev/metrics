package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/konstantin-kukharev/metrics/cmd/agent/report"
	"github.com/konstantin-kukharev/metrics/cmd/agent/service"
	"github.com/konstantin-kukharev/metrics/cmd/agent/settings"
	"github.com/konstantin-kukharev/metrics/cmd/agent/state"
)

func main() {
	conf := settings.New()
	c := &http.Client{}
	s := state.NewMemory()
	r := report.NewRest(c)
	srv := service.NewState(conf, s, r)

	fmt.Printf(
		"runninig agent\r\nreport on %s\r\nreport interval: %s sec.\r\npool interval: %s sec.\r\n",
		conf.GetServerAddress(),
		conf.GetReportInterval(),
		conf.GetPoolInterval(),
	)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

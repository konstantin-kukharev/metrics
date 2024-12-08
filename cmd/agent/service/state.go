package service

import (
	"runtime"
	"time"

	"github.com/konstantin-kukharev/metrics/cmd/agent/report"
	"github.com/konstantin-kukharev/metrics/cmd/agent/settings"
	"github.com/konstantin-kukharev/metrics/cmd/agent/state"
)

type StateService struct {
	app        settings.Application
	nextPool   time.Time
	nextReport time.Time
	m          state.Memory
	r          report.AgentReporter
}

func (s *StateService) Run() error {
	for {
		cTime := time.Now()
		if s.nextPool.Before(cTime) || s.nextPool.Equal(cTime) {
			var mem runtime.MemStats
			s.m.Update(&mem)
			s.nextPool = cTime.Add(s.app.GetPoolInterval())
		}
		if s.nextReport.Before(cTime) || s.nextReport.Equal(cTime) {
			err := s.r.Report(s.app.GetServerAddress(), s.m.Data())
			if err != nil {
				return err
			}
			s.nextReport = cTime.Add(s.app.GetReportInterval())
		}
	}
}

func NewState(
	app settings.Application,
	m state.Memory,
	r report.AgentReporter) *StateService {
	return &StateService{
		app:        app,
		nextPool:   time.Now(),
		nextReport: time.Now(),
		m:          m,
		r:          r,
	}
}

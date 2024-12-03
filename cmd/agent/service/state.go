package service

import (
	"runtime"
	"time"

	"github.com/konstantin-kukharev/metrics/cmd/agent/report"
	"github.com/konstantin-kukharev/metrics/cmd/agent/settings"
	"github.com/konstantin-kukharev/metrics/cmd/agent/state"
)

type stateService struct {
	cfg        settings.Settings
	nextPool   time.Time
	nextReport time.Time
	m          state.Memory
	r          report.AgentReporter
}

func (s *stateService) Run() error {
	for {
		cTime := time.Now()
		if s.nextPool.Before(cTime) || s.nextPool.Equal(cTime) {
			var mem runtime.MemStats
			s.m.Update(&mem)
			s.nextPool = cTime.Add(s.cfg.GetPoolInterval())
		}
		if s.nextReport.Before(cTime) || s.nextReport.Equal(cTime) {
			err := s.r.Report(s.cfg.GetServerAddress(), s.m.Data())
			if err != nil {
				return err
			}
			s.nextReport = cTime.Add(s.cfg.GetReportInterval())
		}
	}
}

func NewState(
	cfg settings.Settings,
	m state.Memory,
	r report.AgentReporter) *stateService {
	return &stateService{
		cfg:        cfg,
		nextPool:   time.Now().Add(cfg.GetPoolInterval()),
		nextReport: time.Now().Add(cfg.GetReportInterval()),
		m:          m,
		r:          r,
	}
}

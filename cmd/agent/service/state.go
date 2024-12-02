package service

import (
	"runtime"
	"time"

	"github.com/konstantin-kukharev/metrics/internal"
)

type stateService struct {
	cfg        internal.AgentSettings
	nextPool   time.Time
	nextReport time.Time
	m          internal.StateMemory
	r          internal.AgentReporter
}

func (s *stateService) Run() {
	for {
		cTime := time.Now()
		if s.nextPool.Before(cTime) || s.nextPool.Equal(cTime) {
			var mem runtime.MemStats
			s.m.Update(&mem)
			s.nextPool = cTime.Add(s.cfg.PoolInterval())
		}
		if s.nextReport.Before(cTime) || s.nextReport.Equal(cTime) {
			s.r.Report(s.cfg.ServerAddress(), s.m.Data())
			s.nextReport = cTime.Add(s.cfg.ReportInterval())
		}
	}
}

func NewState(
	cfg internal.AgentSettings,
	m internal.StateMemory,
	r internal.AgentReporter) *stateService {
	return &stateService{
		cfg:        cfg,
		nextPool:   time.Now(),
		nextReport: time.Now(),
		m:          m,
		r:          r,
	}
}

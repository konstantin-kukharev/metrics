package service

import (
	"runtime"
	"time"

	"github.com/konstantin-kukharev/metrics/internal"
)

type stateService struct {
	poolInterval   time.Duration
	reportInterval time.Duration
	nextPool       time.Time
	nextReport     time.Time
	m              internal.StateMemory
	r              internal.AgentReporter
}

func (s *stateService) Run() {
	for {
		cTime := time.Now()
		if s.nextPool.Before(cTime) || s.nextPool.Equal(cTime) {
			var mem runtime.MemStats
			s.m.Update(&mem)
			s.nextPool = cTime.Add(s.poolInterval)
		}
		if s.nextReport.Before(cTime) || s.nextReport.Equal(cTime) {
			s.r.Report(s.m.Data())
			s.nextReport = cTime.Add(s.reportInterval)
		}
	}
}

func NewState(
	m internal.StateMemory,
	r internal.AgentReporter,
	poolInterval, reportInterval time.Duration) *stateService {
	return &stateService{
		poolInterval:   poolInterval,
		reportInterval: reportInterval,
		nextPool:       time.Now(),
		nextReport:     time.Now(),
		m:              m,
		r:              r,
	}
}

package settings

import (
	"time"
)

type settings struct {
	address        string
	poolInterval   time.Duration
	reportInterval time.Duration
}

func (s *settings) ServerAddress() string {
	return s.address
}

func (s *settings) PoolInterval() time.Duration {
	return s.poolInterval
}

func (s *settings) ReportInterval() time.Duration {
	return s.reportInterval
}

func New() *settings {
	s := new(settings)
	fromFlag(s)

	return s
}

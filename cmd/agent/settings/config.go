package settings

import (
	"time"

	"github.com/konstantin-kukharev/metrics/internal"
)

type Settings interface {
	GetServerAddress() string
	GetReportInterval() time.Duration
	GetPoolInterval() time.Duration
}

type Config struct {
	Address        string
	PoolInterval   time.Duration
	ReportInterval time.Duration
}

func (c *Config) GetServerAddress() string {
	return c.Address
}

func (c *Config) GetPoolInterval() time.Duration {
	return c.PoolInterval
}

func (s *Config) GetReportInterval() time.Duration {
	return s.ReportInterval
}

func New() *Config {
	c := &Config{
		Address:        internal.DefaultServerAddr,
		PoolInterval:   internal.DefaultPoolInterval,
		ReportInterval: internal.DefaultReportInterval,
	}
	fromFlag(c)

	return c
}

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
	PoolInterval   int
	ReportInterval int
}

func (c *Config) GetServerAddress() string {
	return c.Address
}

func (c *Config) GetPoolInterval() time.Duration {
	return time.Duration(c.PoolInterval * int(time.Second))
}

func (c *Config) GetReportInterval() time.Duration {
	return time.Duration(c.ReportInterval * int(time.Second))
}

// Если указана переменная окружения, то используется она.
// Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
// Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
func New() *Config {
	c := &Config{
		Address:        internal.DefaultServerAddr,
		PoolInterval:   internal.DefaultPoolInterval,
		ReportInterval: internal.DefaultReportInterval,
	}
	fromFlag(c)
	fromEnv(c)

	return c
}

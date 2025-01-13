package settings

import (
	"flag"
	"os"
	"strconv"

	"github.com/konstantin-kukharev/metrics/internal"
)

type Config struct {
	Address        string
	PoolInterval   int
	ReportInterval int
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

	return c
}

func (c *Config) WithFlag() *Config {
	flag.StringVar(&c.Address, "a", internal.DefaultServerAddr, "server address")
	flag.IntVar(&c.ReportInterval, "r", internal.DefaultReportInterval, "report interval time duration")
	flag.IntVar(&c.PoolInterval, "p", internal.DefaultPoolInterval, "pool interval time duration")

	flag.Parse()

	return c
}

// ADDRESS отвечает за адрес эндпоинта HTTP-сервера.
// REPORT_INTERVAL позволяет переопределять reportInterval.
// POLL_INTERVAL позволяет переопределять pollInterval.
func (c *Config) WithEnv() *Config {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		c.Address = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if val, err := strconv.Atoi(envReportInterval); err == nil {
			c.ReportInterval = val
		}
	}

	if envPoolInterval := os.Getenv("POLL_INTERVAL"); envPoolInterval != "" {
		if val, err := strconv.Atoi(envPoolInterval); err == nil {
			c.PoolInterval = val
		}
	}

	return c
}

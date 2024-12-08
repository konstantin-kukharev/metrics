package settings

import (
	"flag"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/konstantin-kukharev/metrics/internal"
)

type Application interface {
	GetServerAddress() string
	GetReportInterval() time.Duration
	GetPoolInterval() time.Duration
	Log() Logger
}

type Logger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
}

type Config struct {
	Address        string
	PoolInterval   int
	ReportInterval int
	log            *slog.Logger
	errLog         *slog.Logger
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
	// init application logger
	errHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	})

	c.errLog = slog.New(errHandler)

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})

	c.log = slog.New(logHandler)

	return c
}

func (c *Config) Log() Logger {
	return c.log
}

func (c *Config) WithFlag() *Config {
	flag.StringVar(&c.Address, "a", internal.DefaultServerAddr, "server address")
	flag.IntVar(&c.ReportInterval, "r", 10, "report interval time duration")
	flag.IntVar(&c.PoolInterval, "p", 2, "pool interval time duration")

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

func (c *Config) WithDebug() *Config {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	c.log = logger.With(
		slog.Group("program_info",
			slog.Int("pid", os.Getpid()),
		),
	)

	return c
}

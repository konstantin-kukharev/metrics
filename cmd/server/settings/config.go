package settings

import (
	"flag"
	"log/slog"
	"os"

	"github.com/konstantin-kukharev/metrics/internal"
)

type Application interface {
	GetAddress() string
	Log() Logger
}

type Logger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
}

type Config struct {
	Address string
	log     *slog.Logger
	errLog  *slog.Logger
}

func (c *Config) GetAddress() string {
	return c.Address
}

func NewConfig() *Config {
	c := &Config{Address: internal.DefaultServerAddr}
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
	flag.Parse()

	return c
}

func (c *Config) WithEnv() *Config {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		c.Address = envRunAddr
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

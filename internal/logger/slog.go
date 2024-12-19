package logger

import (
	"log/slog"
	"os"
)

type Slog struct {
	log    *slog.Logger
	errLog *slog.Logger
}

func NewSlog() *Slog {
	l := new(Slog)
	errHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	})

	l.errLog = slog.New(errHandler)

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})

	l.log = slog.New(logHandler)

	return l
}

func (l *Slog) WithDebug(msg string, fields ...any) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	l.log = logger.With(
		slog.Group("program_info",
			slog.Int("pid", os.Getpid()),
		),
	)
}

func (l *Slog) Debug(msg string, fields ...any) {
	l.log.Debug(msg, fields...)
}

func (l *Slog) Info(msg string, fields ...any) {
	l.log.Info(msg, fields...)
}

func (l *Slog) Warn(msg string, fields ...any) {
	l.log.Warn(msg, fields...)
}

func (l *Slog) Error(msg string, fields ...any) {
	l.errLog.Error(msg, fields...)
}

package logger

import (
	"context"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapFieldsKey struct{}

type ZapFields map[string]zap.Field

func (zf ZapFields) Append(fields ...zap.Field) ZapFields {
	zfCopy := make(ZapFields)
	for k, v := range zf {
		zfCopy[k] = v
	}

	for _, f := range fields {
		zfCopy[f.Key] = f
	}

	return zfCopy
}

type settings struct {
	config *zap.Config
	opts   []zap.Option
}

func defaultSettings(level zap.AtomicLevel) *settings {
	config := &zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "@timestamp",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return &settings{
		config: config,
		opts: []zap.Option{
			zap.AddCallerSkip(1),
		},
	}
}

type ZapLogger struct {
	logger *zap.Logger
	level  zap.AtomicLevel
}

// NewZapLogger создает новый логгер Zap.
//
// level - уровень логирования.
//
// Возвращает логгер Zap и ошибку, если возникла ошибка при создании логгера.
func NewZapLogger(level zapcore.Level) (*ZapLogger, error) {
	atomic := zap.NewAtomicLevelAt(level)
	settings := defaultSettings(atomic)

	l, err := settings.config.Build(settings.opts...)
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		logger: l,
		level:  atomic,
	}, nil
}

func (z *ZapLogger) WithContextFields(ctx context.Context, fields ...zap.Field) context.Context {
	ctxFields, _ := ctx.Value(zapFieldsKey{}).(ZapFields)
	if ctxFields == nil {
		ctxFields = make(ZapFields)
	}

	merged := ctxFields.Append(fields...)
	return context.WithValue(ctx, zapFieldsKey{}, merged)
}

func (z *ZapLogger) maskField(f zap.Field) zap.Field {
	if f.Key == "password" {
		return zap.String(f.Key, "******")
	}

	return f
}

func (z *ZapLogger) Sync() {
	_ = z.logger.Sync()
}

func (z *ZapLogger) withCtxFields(ctx context.Context, fields ...zap.Field) []zap.Field {
	fs := make(ZapFields)

	ctxFields, ok := ctx.Value(zapFieldsKey{}).(ZapFields)
	if ok {
		fs = ctxFields
	}

	fs = fs.Append(fields...)

	var maskedFields []zap.Field
	for _, f := range fs {
		maskedFields = append(maskedFields, z.maskField(f))
	}

	return maskedFields
}

func (z *ZapLogger) InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	z.logger.Info(msg, z.withCtxFields(ctx, fields...)...)
}

func (z *ZapLogger) DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	z.logger.Debug(msg, z.withCtxFields(ctx, fields...)...)
}

func (z *ZapLogger) WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	z.logger.Warn(msg, z.withCtxFields(ctx, fields...)...)
}

func (z *ZapLogger) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	z.logger.Error(msg, z.withCtxFields(ctx, fields...)...)
}

func (z *ZapLogger) FatalCtx(ctx context.Context, msg string, fields ...zap.Field) {
	z.logger.Fatal(msg, z.withCtxFields(ctx, fields...)...)
}

func (z *ZapLogger) PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	z.logger.Panic(msg, z.withCtxFields(ctx, fields...)...)
}

func (z *ZapLogger) SetLevel(level zapcore.Level) {
	z.level.SetLevel(level)
}

func (z *ZapLogger) Std() *log.Logger {
	return zap.NewStdLog(z.logger)
}

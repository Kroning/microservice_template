package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global *zap.Logger
	level  = zap.NewAtomicLevelAt(zap.DebugLevel)
)

func init() {
	SetLogger(New(level))
}

func New(lvl zapcore.LevelEnabler, options ...zap.Option) *zap.Logger {
	if lvl == nil {
		lvl = level
	}
	sink := zapcore.AddSync(os.Stdout)
	options = append(options, zap.AddCallerSkip(1), zap.AddCaller())
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				TimeKey:        "ts",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			sink,
			lvl,
		),
		options...,
	)
}

func SetLevel(l string) {
	var zapLevel zapcore.Level
	if err := zapLevel.Set(l); err != nil {
		Warn(context.Background(), "failed to parse, setting to WARN", zap.Reflect("level", level), zap.Error(err))
		zapLevel = zapcore.WarnLevel
	}

	level.SetLevel(zapLevel)
}

func Logger() *zap.Logger {
	return global
}

func SetLogger(l *zap.Logger) {
	global = l
}

func Debug(ctx context.Context, message string, kvs ...zap.Field) {
	global.Debug(message, kvs...)
}

func Info(ctx context.Context, message string, kvs ...zap.Field) {
	global.Info(message, kvs...)
}

func Warn(ctx context.Context, message string, kvs ...zap.Field) {
	global.Warn(message, kvs...)
}

func Error(ctx context.Context, message string, kvs ...zap.Field) {
	global.Error(message, kvs...)
}

func Fatal(ctx context.Context, message string, kvs ...zap.Field) {
	global.Fatal(message, kvs...)
}

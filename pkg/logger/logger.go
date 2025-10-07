package logger

import (
	"context"

	"go.uber.org/zap"
)

var DefaultLogger Logger = &zapLogger{logger: zap.L().Sugar()}

// Logger определяет интерфейс для логирования
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})

	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})

	// With добавляет поля к логгеру
	With(args ...interface{}) Logger
	WithTrace(ctx context.Context) Logger

	ToContext(ctx context.Context) context.Context
}

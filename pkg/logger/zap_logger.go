package logger

import (
	"context"

	"github.com/curserio/chrono-api/pkg/tracer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger создает новый логгер на основе zap
func NewZapLogger(isDevelopment bool) (*zap.SugaredLogger, error) {
	var logger *zap.Logger
	var err error

	if isDevelopment {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.DisableStacktrace = true
		logger, err = config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		config := zap.NewProductionConfig()
		config.DisableStacktrace = true
		logger, err = config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	}

	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

type zapLogger struct {
	logger    *zap.SugaredLogger
	withTrace bool
}

func AdaptZap(logger *zap.SugaredLogger) (Logger, error) {
	adapt := &zapLogger{logger: logger}
	DefaultLogger = adapt
	return adapt, nil
}

func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) Debugf(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

func (l *zapLogger) Infof(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

func (l *zapLogger) Warnf(msg string, args ...interface{}) {
	l.logger.Warnf(msg, args...)
}

func (l *zapLogger) Errorf(msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

func (l *zapLogger) Fatalf(msg string, args ...interface{}) {
	l.logger.Fatalf(msg, args...)
}

func (l *zapLogger) With(args ...interface{}) Logger {
	newLogger := l.logger.With(args...)
	return &zapLogger{logger: newLogger, withTrace: l.withTrace}
}

func (l *zapLogger) WithTrace(ctx context.Context) Logger {
	if l.withTrace {
		return l
	}

	traceID := tracer.TraceIDFromContext(ctx)
	if traceID == "" {
		return l
	}

	newLogger := l.logger.With("trace_id", traceID)
	return &zapLogger{logger: newLogger, withTrace: true}
}

func (l *zapLogger) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

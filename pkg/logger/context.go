package logger

import (
	"context"
)

type loggerKey struct{}

func FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return DefaultLogger
	}

	log := DefaultLogger

	if l, ok := ctx.Value(loggerKey{}).(Logger); ok {
		log = l
	}

	return log.WithTrace(ctx)
}

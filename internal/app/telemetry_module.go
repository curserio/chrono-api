package app

import (
	"context"

	"github.com/curserio/chrono-api/config"
	"github.com/curserio/chrono-api/pkg/tracer"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

var TelemetryModule = fx.Options(
	fx.Provide(
		func(lc fx.Lifecycle, cfg *config.Config) (*sdktrace.TracerProvider, error) {
			tp, err := tracer.NewTracerProvider(cfg.App.Name)
			if err != nil {
				return nil, err
			}

			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					return tp.Shutdown(ctx)
				},
			})

			return tp, nil
		},
	),
	fx.Invoke(func(*sdktrace.TracerProvider) {}),
)

package main

import (
	"github.com/curserio/chrono-api/config"
	"github.com/curserio/chrono-api/internal/app"
	"github.com/curserio/chrono-api/internal/infrastructure/http/server"
	"github.com/curserio/chrono-api/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// TODO
// TODOS: хранение языка
// TODOS:
func main() {
	fx.New(
		// Logger used by Fx
		fx.WithLogger(func(log *zap.SugaredLogger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Desugar()}
		}),

		fx.Provide(
			// Config
			config.NewConfig,

			// Logger
			func(cfg *config.Config) (*zap.SugaredLogger, error) {
				return logger.NewZapLogger(cfg.App.DevMode)
			},
			logger.AdaptZap,

			// Server
			func(l logger.Logger, cfg *config.Config) []func(*server.Server) {
				return []func(*server.Server){
					server.WithLogger(l),
					server.WithDefaultLanguage(cfg.App.DefaultLanguage),
				}
			},

			// Server constructor
			func(lc fx.Lifecycle, cfg *config.Config, opts []func(*server.Server)) *server.Server {
				return server.New(cfg.Server.Addr, cfg.App.Name, lc, opts...)
			},
		),

		app.RepositoryModule,
		app.UsecaseModule,
		app.HandlerModule,
		app.TelemetryModule,
	).Run()
}

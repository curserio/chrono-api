package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/curserio/chrono-api/internal/middleware"
	"github.com/curserio/chrono-api/pkg/logger"
	"github.com/curserio/chrono-api/pkg/validator"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/fx"
)

type Server struct {
	echo *echo.Echo

	addr            string
	defaultLanguage string
	logger          logger.Logger
}

func New(addr, serviceName string, lc fx.Lifecycle, opts ...func(*Server)) *Server {
	server := &Server{addr: addr, defaultLanguage: defaultLanguage}

	for _, opt := range opts {
		opt(server)
	}

	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = server.errorHandler

	v := validator.NewValidator()
	v.RegisterCustomValidations() // регистрируем кастомные правила
	e.Validator = v

	// tracing
	e.Use(otelecho.Middleware(serviceName))
	// logging
	if server.logger != nil {
		e.Use(middleware.RequestLogger(server.logger))
	}
	// language
	e.Use(middleware.I18nMiddleware(server.defaultLanguage))

	server.echo = e

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if server.logger != nil {
				server.logger.Info(fmt.Sprintf("server starting at %s", server.addr))
			}

			go func() {
				err := server.Start()
				if err != nil && !errors.Is(http.ErrServerClosed, err) {
					if server.logger != nil {
						server.logger.Fatal("server failed to start", "error", err)
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if server.logger != nil {
				server.logger.Info("server shutting down")
			}
			return server.Shutdown(ctx)
		},
	})

	return server
}

func (s *Server) NewGroup(prefix string, m ...echo.MiddlewareFunc) *echo.Group {
	return s.echo.Group(prefix, m...)
}

func (s *Server) Start() error {
	return s.echo.Start(s.addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

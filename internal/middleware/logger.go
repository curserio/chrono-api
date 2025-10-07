package middleware

import (
	"time"

	"github.com/curserio/chrono-api/pkg/logger"
	"github.com/labstack/echo/v4"
)

// RequestLogger создает middleware для логирования HTTP запросов
func RequestLogger(baseLogger logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			l := baseLogger.WithTrace(req.Context()).With(
				"method", req.Method,
				"uri", req.RequestURI,
				"remote_addr", req.RemoteAddr,
			)

			ctx := l.ToContext(req.Context())
			c.SetRequest(req.WithContext(ctx))

			// Логируем начало запроса
			l.Info("incoming request")

			err := next(c)

			fields := []interface{}{"duration", time.Since(start).String()}

			// при ошибках статус извлекаем из ошибки в errorHandler
			if err == nil {
				fields = append(fields, []interface{}{"status", res.Status}...)
			}

			l = l.With(fields...)

			ctx = l.ToContext(req.Context())
			c.SetRequest(req.WithContext(ctx))

			// Логируем завершение запроса
			l.Info("request completed")

			return err
		}
	}
}

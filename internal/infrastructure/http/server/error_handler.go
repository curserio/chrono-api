package server

import (
	"net/http"

	apiErrors "github.com/curserio/chrono-api/internal/errors"
	"github.com/curserio/chrono-api/pkg/logger"
	"github.com/labstack/echo/v4"
)

func (s *Server) errorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var (
		statusCode = echo.ErrInternalServerError.Code
		apiErr     = apiErrors.ApiError{}
	)

	httpErr, ok := apiErrors.Unwrap(err)
	if ok {
		statusCode = httpErr.Code
		apiErr.Error = httpErr.Message
	} else {
		statusCode = http.StatusInternalServerError
		apiErr.Error = err.Error()
	}

	log := logger.FromContext(c.Request().Context())

	log.Error("handling request error",
		"error", err,
		"error_message", apiErr.Error,
		"status", statusCode,
	)

	if err = c.JSON(statusCode, apiErr); err != nil {
		log.Error("send error response", "error", err)
	}
}

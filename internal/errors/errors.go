package errors

import (
	"errors"
	"net/http"
	"time"
)

var (
	ErrNotFound = errors.New("resource not found")

	ErrBookingStatusInvalid   = errors.New("invalid booking status")
	ErrScheduleTypeInvalid    = errors.New("invalid schedule type")
	ErrEndTimeBeforeStartTime = errors.New("end time is before start time")
)

type HTTPError struct {
	Code       int
	Message    string
	InnerError error
	Timestamp  time.Time
}

func NewHTTPError(code int, message string, inner error) *HTTPError {
	return &HTTPError{
		Code:       code,
		Message:    message,
		InnerError: inner,
		Timestamp:  time.Now(),
	}
}

func (e *HTTPError) Error() string {
	if e.InnerError != nil {
		return e.InnerError.Error()
	}
	return e.Message
}

func GetCodeAndMessage(err error) (int, string) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.Code, httpErr.Message
	}
	return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
}

func Unwrap(err error) (*HTTPError, bool) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr, true
	}
	return nil, false
}

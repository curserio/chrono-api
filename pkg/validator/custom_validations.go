package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// registerTimeValidations регистрирует кастомные правила валидации для времени
func (cv *CustomValidator) registerTimeValidations() error {
	return cv.validator.RegisterValidation("futuretime", validateFutureTime)
}

// validateFutureTime проверяет, что время находится в будущем
func validateFutureTime(fl validator.FieldLevel) bool {
	timeVal, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return timeVal.After(time.Now())
}

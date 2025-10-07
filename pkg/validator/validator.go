package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator представляет собой интерфейс для валидации
type Validator interface {
	Validate(i interface{}) error
}

// CustomValidator реализует интерфейс echo.Validator
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator создает новый экземпляр валидатора
func NewValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate выполняет валидацию структуры
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// RegisterCustomValidations регистрирует кастомные правила валидации
func (cv *CustomValidator) RegisterCustomValidations() {
	err := cv.registerTimeValidations()
	if err != nil {
		panic(err)
	}
}

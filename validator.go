package sw

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			type FieldError struct {
				Name    string `json:"name"`
				Message string `json:"message"`
			}

			var fieldErrors []FieldError

			for _, e := range validationErrors {
				field := toSnakeCase(e.Field())
				var msg string

				switch e.Tag() {
				case "required":
					msg = fmt.Sprintf("%s is required", field)
				case "e164":
					msg = fmt.Sprintf("%s must be a valid phone number (e.g. +628123456789)", field)
				case "numeric":
					msg = fmt.Sprintf("%s must be numeric", field)
				case "oneof":
					msg = fmt.Sprintf("%s must be one of [%s]", field, e.Param())
				default:
					msg = fmt.Sprintf("%s is invalid", field)
				}

				fieldErrors = append(fieldErrors, FieldError{
					Name:    field,
					Message: msg,
				})
			}

			return echo.NewHTTPError(http.StatusUnprocessableEntity, map[string]interface{}{
				"message": "validation failed",
				"errors":  fieldErrors,
			})
		}
		return err
	}
	return nil
}

// toSnakeCase turns "PhoneNumber" → "phone_number"
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_', r+('a'-'A'))
		} else {
			result = append(result, r)
		}
	}
	return strings.ToLower(string(result))
}

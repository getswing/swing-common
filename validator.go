package sw

import (
	"fmt"
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
			var messages []string
			for _, e := range validationErrors {
				switch e.Tag() {
				case "required":
					messages = append(messages, fmt.Sprintf("%s is required", toSnakeCase(e.Field())))
				case "e164":
					messages = append(messages, fmt.Sprintf("%s must be a valid phone number (e.g. +628123456789)", toSnakeCase(e.Field())))
				case "numeric":
					messages = append(messages, fmt.Sprintf("%s must be numeric", toSnakeCase(e.Field())))
				case "oneof":
					messages = append(messages, fmt.Sprintf("%s must be one of [%s]", toSnakeCase(e.Field()), e.Param()))
				default:
					messages = append(messages, fmt.Sprintf("%s is invalid", toSnakeCase(e.Field())))
				}
			}
			return echo.NewHTTPError(400, strings.Join(messages, ", "))
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

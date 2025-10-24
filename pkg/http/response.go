package common_http

import (
	"github.com/labstack/echo/v4"
)

type SuccessResponseDTO struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponseDTO struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func SuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, SuccessResponseDTO{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, ErrorResponseDTO{
		Status:  "error",
		Message: message,
		Code:    statusCode,
	})
}

package sw

import (
	"github.com/labstack/echo/v4"
)

type SuccessResponseDTO struct {
	Status    string      `json:"status"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	RequestId string      `json:"request_id"`
}

type ErrorResponseDTO struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	RequestId string `json:"request_id"`
	ErrorCode string `json:"error_code,omitempty"`
}

func SuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	response := SuccessResponseDTO{
		Status:    "success",
		Message:   message,
		Data:      data,
		RequestId: requestId,
	}

	c.Logger().Info("SuccessResponse: %v", response)
	return c.JSON(statusCode, response)
}

func ErrorResponse(c echo.Context, statusCode int, message string, errorCode string) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	response := ErrorResponseDTO{
		Status:    "error",
		Message:   message,
		Code:      statusCode,
		RequestId: requestId,
		ErrorCode: errorCode,
	}
	c.Logger().Info("ErrorResponse: %v", response)
	return c.JSON(statusCode, response)
}

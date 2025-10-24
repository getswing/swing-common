package common_http

import (
	common_logger "getswing.app/player-service/pkg/logger"
	"github.com/labstack/echo/v4"
)

func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID := c.Request().Header.Get("X-Request-ID")
		ctx := c.Request().Context()
		ctx = common_logger.WithRequestIDFromHeader(ctx, reqID)
		c.SetRequest(c.Request().WithContext(ctx))
		c.Set("ctx", ctx)
		c.Response().Header().Set("X-Request-ID", common_logger.GetRequestID(ctx))
		return next(c)
	}
}

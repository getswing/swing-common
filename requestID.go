package sw

import (
	sw "github.com/getswing/swing-common"
	"github.com/labstack/echo/v4"
)

func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID := c.Request().Header.Get("X-Request-ID")
		ctx := c.Request().Context()
		ctx = sw.WithRequestIDFromHeader(ctx, reqID)
		c.SetRequest(c.Request().WithContext(ctx))
		c.Set("ctx", ctx)
		c.Response().Header().Set("X-Request-ID", sw.GetRequestID(ctx))
		return next(c)
	}
}

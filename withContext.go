package sw

import (
	"context"

	"github.com/labstack/echo/v4"
)

func WithContext(ctx context.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)
			return next(c)
		}
	}
}

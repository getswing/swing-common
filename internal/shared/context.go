package shared

import (
	"getswing.app/player-service/internal/config"

	"github.com/labstack/echo/v4"
)

const ContextConfigKey = "app.config"

// WithConfig injects a pointer to config.Config into Echo context for every request.
func WithConfig(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(ContextConfigKey, cfg)
			return next(c)
		}
	}
}

// GetConfig fetches the shared config from Echo context.
func GetConfig(c echo.Context) *config.Config {
	if v, ok := c.Get(ContextConfigKey).(*config.Config); ok {
		return v
	}
	return nil
}

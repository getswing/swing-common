package router

import (
	"database/sql"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"getswing.app/player-service/internal/config"
	"getswing.app/player-service/internal/db"
	"getswing.app/player-service/internal/handlers"
	"getswing.app/player-service/internal/mq"
)

// Register wires all public and protected routes and middlewares.
func Register(
	e *echo.Echo,
	cfg config.Config,
	sqlDB *sql.DB,
	amqpCh *amqp.Channel,
	authHandler *handlers.AuthHandler,
) {
	// Public routes
	e.GET("/health", func(c echo.Context) error {
		if err := db.Ping(c.Request().Context(), sqlDB); err != nil {
			return c.JSON(http.StatusServiceUnavailable, echo.Map{"status": "unhealthy", "db": err.Error()})
		}
		if err := mq.Ping(amqpCh); err != nil {
			return c.JSON(http.StatusServiceUnavailable, echo.Map{"status": "unhealthy", "mq": err.Error()})
		}
		return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
	})

	// Auth endpoints
	e.POST("/auth/register", authHandler.Register)
	e.POST("/auth/login", authHandler.Login)

	// Protected API routes with JWT
	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte(cfg.JWTSecret),
		SigningMethod: "HS256",
		TokenLookup:   "header:Authorization",
		AuthScheme:    "Bearer",
	}))

}

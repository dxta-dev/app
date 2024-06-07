package middleware

import (
	"github.com/dxta-dev/app/internal/config"

	"context"

	"github.com/labstack/echo/v4"
)

const ConfigContextKey string = "app_config"

func ConfigMiddleware(config *config.Config) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			ctx := r.Context()
			ctx = context.WithValue(ctx, ConfigContextKey, config)
			c.SetRequest(r.WithContext(ctx))
			return next(c)
		}
	}
}

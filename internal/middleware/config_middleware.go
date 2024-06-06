package middleware

import (
	"github.com/dxta-dev/app/internal/config"

	"context"

	"github.com/labstack/echo/v4"
)

type configContextKey string

const ConfigContext configContextKey = "app_config"

func WithConfigContext(parent context.Context, config *config.Config) context.Context {
	return context.WithValue(parent, ConfigContext, config)
}

func ConfigMiddleware(config *config.Config) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			ctx := r.Context()
			ctx = WithConfigContext(ctx, config)
			c.SetRequest(r.WithContext(ctx))
			return next(c)
		}
	}
}

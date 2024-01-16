package middlewares

import (
	"context"
	"dxta-dev/app/internal/utils"

	"github.com/labstack/echo/v4"
)

type configContextKey string

const ConfigContext configContextKey = "app_config"

func ConfigMiddleware(config *utils.Config) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			ctx := r.Context()
			ctx = context.WithValue(ctx, ConfigContext, config)
			c.SetRequest(r.WithContext(ctx))
			return next(c)
		}
	}
}

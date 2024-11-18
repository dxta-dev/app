package middleware

import (
	"github.com/dxta-dev/app/internal/otel"
	"github.com/dxta-dev/app/internal/web/data"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
)

const StoreContextKey string = "store"

func StoreMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			ctx := r.Context()
			tenantDatabaseURL := ctx.Value(TenantDatabaseURLContext).(string)

			driverName := otel.GetDriverName()

			store := &data.Store{
				DbUrl:      tenantDatabaseURL,
				DriverName: driverName,
				Context:    ctx,
			}

			ctx = context.WithValue(ctx, StoreContextKey, store)
			c.SetRequest(r.WithContext(ctx))
			return next(c)
		}
	}
}

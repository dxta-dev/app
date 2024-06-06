package middleware

import (
	"github.com/labstack/echo/v4"
)

const TenantDatabaseUrlKey string = "tenant_database_url"

func TenantMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

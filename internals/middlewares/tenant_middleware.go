package middlewares

import (
	"github.com/labstack/echo/v4"
	"strings"
)

func TenantMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		host := c.Request().Host
		parts := strings.Split(host, ".")

		if len(parts) > 2 {
			subdomain := parts[0]
			c.Set("subdomain", subdomain)
			c.Set("is_root", false)
		} else {
			c.Set("subdomain", "root")
			c.Set("is_root", true)
		}

		return next(c)
	}
}

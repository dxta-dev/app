package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

type TenantInfo struct {
	Tenants []string `json:"tenants"`
}

var tenantsMap = make(map[string]bool)

const SubdomainContextKey = string("subdomain")
const IsRootContextKey = string("is_root")
const TenantContextKey = string("tenant")

func LoadTenants() error {
	resp, err := http.Get(os.Getenv("OSS_TENANTS_ENDPOINT"))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var tenants TenantInfo

	if err := json.Unmarshal(body, &tenants); err != nil {
		return err
	}

	for _, tenant := range tenants.Tenants {
		tenantsMap[tenant] = true
	}

	return nil
}

func TenantMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		hostName := c.Request().Host
		parts := strings.Split(hostName, ".")

		if len(parts) <= 2 {
			ctx = context.WithValue(ctx, SubdomainContextKey, "root")
			ctx = context.WithValue(ctx, IsRootContextKey, true)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
		tenant := parts[0]

		if _, exists := tenantsMap[tenant]; !exists {
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("http://%s", strings.Join(parts[1:], ".")))
		}

		ctx = context.WithValue(ctx, SubdomainContextKey, tenant)
		ctx = context.WithValue(ctx, TenantContextKey, tenant)
		ctx = context.WithValue(ctx, IsRootContextKey, false)

		// Use the new context in the request
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

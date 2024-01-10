package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type TenantInfo struct {
	Tenants []string `json:"tenants"`
}

type TenantMap map[string]bool

var tenantsMap = make(TenantMap)

type subdomainContextKey string
type isRootContextKey string
type tenantContextKey string

const SubdomainContext subdomainContextKey = "subdomain"
const IsRootContext isRootContextKey = "is_root"
const TenantContext tenantContextKey = "tenant"

func LoadTenantsDummy(tenantsMapArg TenantMap) {
	tenantsMap = tenantsMapArg
}

func LoadTenantsFromAPI(ossTenantsEndpoint string) error {
	resp, err := http.Get(ossTenantsEndpoint)
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
			ctx = context.WithValue(ctx, SubdomainContext, "root")
			ctx = context.WithValue(ctx, IsRootContext, true)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
		tenant := parts[0]

		if _, exists := tenantsMap[tenant]; !exists {
			// TODO: http/https
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("http://%s", strings.Join(parts[1:], ".")))
		}

		ctx = context.WithValue(ctx, SubdomainContext, tenant)
		ctx = context.WithValue(ctx, TenantContext, tenant)
		ctx = context.WithValue(ctx, IsRootContext, false)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

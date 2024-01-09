package middlewares

import (
	"context"
	"database/sql"
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

var tenantsMap = make(map[string]*sql.DB)

const TenantDatabaseContext = "Tenant DB"

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
		db, err := sql.Open("libsql", fmt.Sprintf("libsql://%s-dxta.turso.io?authToken=%s", tenant, os.Getenv("DATABASE_AUTH_TOKEN")))
		if err != nil {
			return err
		}

		// TODO Do we need manual clean up?
		defer db.Close()
		tenantsMap[tenant] = db
	}

	return nil
}

func TenantMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		hostName := c.Request().Host
		parts := strings.Split(hostName, ".")

		if len(parts) <= 2 {
			ctx = context.WithValue(ctx, "subdomain", "root")
			ctx = context.WithValue(ctx, "is_root", true)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
		tenant := parts[0]
		db, tenantExists := tenantsMap[tenant]
		if !tenantExists {
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("http://%s", strings.Join(parts[1:], ".")))
		}

		ctx = context.WithValue(ctx, "subdomain", tenant)
		ctx = context.WithValue(ctx, "is_root", false)
		ctx = context.WithValue(ctx, TenantDatabaseContext, db)

		// Use the new context in the request
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

package middlewares

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

type TenantInfo struct {
	Name string `json:"tenant"`
}

type TenantsConfig struct {
	tenants []TenantInfo
}

// TODO: just key-value store map (tenant-key)->(db-connection)
var tenantsConfig TenantsConfig

func LoadTenants() error {
	resp, err := http.Get(os.Getenv("OSS_TENANTS_ENPDOINT"))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &tenantsConfig.tenants); err != nil {
		return err
	}

	// TODO: open db connections ?
	return nil
}

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

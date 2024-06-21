package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dxta-dev/app/internal/config"
	"github.com/labstack/echo/v4"
)

type hostInfo struct {
	ProtocolScheme string
	Subdomain      string
	IsRoot         bool
}

func getHostInfo(r *http.Request) hostInfo {
	tls := r.TLS
	hostName := r.Host
	parts := strings.Split(hostName, ".")
	hostProtocolScheme := "https"
	if tls == nil {
		hostProtocolScheme = "http"
	}

	subdomain, isRoot := parts[0], false

	if len(parts) <= 2 {
		subdomain = ""
		isRoot = true
	}

	return hostInfo{
		ProtocolScheme: hostProtocolScheme,
		Subdomain:      subdomain,
		IsRoot:         isRoot,
	}
}

const TenantDatabaseUrlKey string = "tenant_database_url"

func TenantMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		conf := ctx.Value(ConfigContextKey).(*config.Config)

		if conf == nil {
			return errors.New("config not found in context")
		}

		hostInfo := getHostInfo(c.Request())

		if conf.IsSuperDatabaseEnabled {


			if hostInfo.Subdomain == "" {

			}

		}

		return next(c)
	}
}

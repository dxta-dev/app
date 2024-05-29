package middleware

import (
	"github.com/dxta-dev/app/internal/otel"
	"github.com/dxta-dev/app/internal/util"

	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type TenantDbUrlMap map[string]string

type subdomainContextKey string
type isRootContextKey string
type tenantContextKey string

const SubdomainContext subdomainContextKey = "subdomain"
const IsRootContext isRootContextKey = "is_root"
const TenantDatabaseURLContext tenantContextKey = "tenant_db_url"
const TenantDatabasesGlobalContext string = "tenant_db_map"

// TODO(scalability?): change from const map fetch to cache per tenant?
func getTenantToDatabaseURLMap(ctx context.Context, superDatabaseUrl string) (TenantDbUrlMap, error) {
	tenantToDatabaseURLMap := make(TenantDbUrlMap)

db, err := sql.Open(otel.GetDriverName(), superDatabaseUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	query := `
		SELECT
			subdomain,
			db_url
		FROM tenants
	`

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var tenantName string
		var tenantDatabaseURL string

		if err := rows.Scan(&tenantName, &tenantDatabaseURL); err != nil {
			log.Fatal((err))
		}

		tenantToDatabaseURLMap[tenantName] = tenantDatabaseURL
	}

	return tenantToDatabaseURLMap, nil
}

func getTenantDatabaseURL(ctx context.Context, config *util.Config, tenantKey string) (string, bool, error) {

	if !config.ShouldUseSuperDatabase {
		configTenant, configContainsTenant := config.Tenants[tenantKey]

		if !configContainsTenant {
			return "", false, nil
		}

		return *configTenant.DatabaseUrl, true, nil
	}

	tenantsToDatabaseURLMap, err := getTenantToDatabaseURLMap(ctx, *config.SuperDatabaseUrl)

	if err != nil {
		return "", false, err
	}

	tenantDatabaseURL, databaseContainsTenant := tenantsToDatabaseURLMap[tenantKey]

	if !databaseContainsTenant {
		return "", false, nil
	}

	return fmt.Sprintf(*config.TenantDatabaseUrlTemplate, tenantDatabaseURL), true, nil
}

func TenantMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		config := ctx.Value(ConfigContext).(*util.Config)
		tls := c.Request().TLS
		hostName := c.Request().Host
		parts := strings.Split(hostName, ".")
		hostProtocolScheme := "https"
		if tls == nil {
			hostProtocolScheme = "http"
		}

		subdomain, isRoot := parts[0], false

		if len(parts) <= 2 {
			subdomain = "root"
			isRoot = true
		}

		ctx = context.WithValue(ctx, SubdomainContext, subdomain)
		ctx = context.WithValue(ctx, IsRootContext, isRoot)

		if !config.IsMultiTenant && len(config.Tenants) == 1 {
			for _, v := range config.Tenants {
				ctx = context.WithValue(ctx, TenantDatabaseURLContext, *v.DatabaseUrl)
			}

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}

		if isRoot {
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}

		tenantDatabaseUrl, tenantDatabaseUrlExists, err := getTenantDatabaseURL(ctx, config, subdomain)

		if err != nil {
			log.Panicln("Error getting tenant database URL", err)
		}

		if !tenantDatabaseUrlExists {
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s://%s/oss", hostProtocolScheme, strings.Join(parts[1:], ".")))
		}

		ctx = context.WithValue(ctx, TenantDatabaseURLContext, tenantDatabaseUrl)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

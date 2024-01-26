package middlewares

import (
	"context"
	"database/sql"
	"dxta-dev/app/internal/utils"
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
func getTenantToDatabaseURLMap(superDatabaseUrl string) (TenantDbUrlMap, error) {
	tenantToDatabaseURLMap := make(TenantDbUrlMap)

	db, err := sql.Open("libsql", superDatabaseUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, err
	}

	query := `
		SELECT
			name,
			db_url
		FROM tenants
	`

	rows, err := db.Query(query)

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

func getTenantDatabaseURL(config *utils.Config, tenantKey string) (string, bool, error) {

	if !config.ShouldUseSuperDatabase {
		configTenant, configContainsTenant := config.Tenants[tenantKey]

		if !configContainsTenant {
			return "", false, nil
		}

		return *configTenant.DatabaseUrl, true, nil
	}

	tenantsToDatabaseURLMap, err := getTenantToDatabaseURLMap(*config.SuperDatabaseUrl)

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
		config := ctx.Value(ConfigContext).(*utils.Config)
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

		tenantDatabaseUrl, tenantDatabaseUrlExists, err := getTenantDatabaseURL(config, subdomain)

		if err != nil {
			fmt.Println("Error multi_tenant_middleware.go: TODO(error-handling) - log or something when super database fails")
			// Ideas: https://echo.labstack.com/docs/error-handling
			return echo.ErrInternalServerError
		}

		if !tenantDatabaseUrlExists {
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s://%s/oss", hostProtocolScheme, strings.Join(parts[1:], ".")))
		}

		ctx = context.WithValue(ctx, TenantDatabaseURLContext, tenantDatabaseUrl)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

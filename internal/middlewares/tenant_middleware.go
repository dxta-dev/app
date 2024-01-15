package middlewares

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
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

// TODO(scalability?): change from const map fetch to LRU cache filling
func getTenantToDatabaseURLMap() (TenantDbUrlMap, error) {

	tenantToDatabaseURLMap := make(TenantDbUrlMap)

	db, err := sql.Open("libsql", os.Getenv("SUPER_DATABASE_URL"))

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

func lazyloadTenantToDatabaseURLMap(c echo.Context) (TenantDbUrlMap, error) {
	tenantToDatabaseURLMap, exists := c.Get(TenantDatabasesGlobalContext).(TenantDbUrlMap)

	if exists {
		return tenantToDatabaseURLMap, nil
	}
	// TODO: exists is never true

	tenantToDatabaseURLMap, err := getTenantToDatabaseURLMap()
	if err != nil {
		return nil, err
	}
	c.Set(TenantDatabasesGlobalContext, tenantToDatabaseURLMap)

	return tenantToDatabaseURLMap, nil
}

func TenantMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		tls := c.Request().TLS
		hostName := c.Request().Host
		parts := strings.Split(hostName, ".")
		hostProtocolScheme := "https"
		if tls == nil {
			hostProtocolScheme = "http"
		}

		if len(parts) <= 2 {
			ctx = context.WithValue(ctx, SubdomainContext, "root")
			ctx = context.WithValue(ctx, IsRootContext, true)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}

		tenant := parts[0]

		tenantToDatabaseURLMap, err := lazyloadTenantToDatabaseURLMap(c)
		if err != nil {
			// TODO(error-handling): log or something
			// Ideas: https://echo.labstack.com/docs/error-handling
			return echo.ErrInternalServerError
		}

		tenantDatabaseUrl, tenantDatabaseUrlExists := tenantToDatabaseURLMap[tenant]
		if !tenantDatabaseUrlExists {
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s://%s/oss", hostProtocolScheme, strings.Join(parts[1:], ".")))
		}

		ctx = context.WithValue(ctx, SubdomainContext, tenant)
		ctx = context.WithValue(ctx, TenantDatabaseURLContext, fmt.Sprintf("%s?authToken=%s", tenantDatabaseUrl, os.Getenv("TENANT_DATABASE_AUTH_TOKEN")))
		ctx = context.WithValue(ctx, IsRootContext, false)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

package middleware

import (
	"github.com/dxta-dev/app/internal/config"
	"github.com/dxta-dev/app/internal/otel"

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


func getTenantsFromSuperDatabase(ctx context.Context, superDatabaseUrl string) ([]config.Tenant, error) {
}


func TenantMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

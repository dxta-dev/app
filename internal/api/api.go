package api

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/dxta-dev/app/internal/api/data"
	"github.com/dxta-dev/app/internal/otel"
	"github.com/labstack/echo/v4"

	_ "github.com/libsql/libsql-client-go/libsql"
)

type APIState struct {
	DB     data.DB
	Org    string
	Repo   string
	TeamId *int64
}

var tenantRepoCache sync.Map

func getCachedTenantRepo(ctx context.Context, org, repo string) (data.TenantRepo, error) {
	cacheKey := strings.ToLower(org + "/" + repo)

	if cachedUrl, ok := tenantRepoCache.Load(cacheKey); ok {
		return cachedUrl.(data.TenantRepo), nil
	}
	driverName := otel.GetDriverName()

	reposDB, err := sql.Open(driverName, os.Getenv("SUPER_DATABASE_URL")+"?authToken="+os.Getenv("DXTA_DEV_GROUP_TOKEN"))
	if err != nil {
		return data.TenantRepo{}, err
	}
	defer reposDB.Close()

	tenantRepo, err := data.GetTenantRepo(ctx, reposDB, org, repo)
	if err != nil {
		return data.TenantRepo{}, err
	}

	tenantRepoCache.Store(cacheKey, tenantRepo)

	return tenantRepo, nil
}

func NewAPIState(c echo.Context) (APIState, error) {

	ctx := c.Request().Context()

	org := c.Param("org")
	repo := c.Param("repo")

	if org == "" || repo == "" {
		return APIState{}, echo.NewHTTPError(http.StatusBadRequest, "org and repo are required")
	}

	tenantRepo, err := getCachedTenantRepo(ctx, org, repo)
	if err != nil {
		return APIState{}, err
	}

	db, err := data.NewDB(ctx, tenantRepo)

	if err != nil {
		return APIState{}, err
	}

	team := c.QueryParam("team")

	var teamInt *int64

	if t, err := strconv.ParseInt(team, 10, 64); err == nil && t > 0 {
		teamInt = &t
	}

	return APIState{
		DB:     db,
		Org:    tenantRepo.Organization,
		Repo:   tenantRepo.Repository,
		TeamId: teamInt,
	}, nil
}

func GetReposDB() (*sql.DB, error) {
	driverName := otel.GetDriverName()

	reposDB, err := sql.Open(driverName, os.Getenv("SUPER_DATABASE_URL")+"?authToken="+os.Getenv("DXTA_DEV_GROUP_TOKEN"))
	if err != nil {
		return nil, err
	}

	return reposDB, nil
}

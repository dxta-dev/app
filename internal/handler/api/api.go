package api

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/otel"
	"github.com/labstack/echo/v4"
)

type APIState struct {
	DB     *sql.DB
	org    string
	repo   string
	teamId *int64
}

var reposDBCache sync.Map

func getCachedDbUrl(ctx context.Context, org, repo string) (string, error) {
	cacheKey := org + "/" + repo

	if cachedUrl, ok := reposDBCache.Load(cacheKey); ok {
		return cachedUrl.(string), nil
	}
	driverName := otel.GetDriverName()

	reposDB, err := sql.Open(driverName, os.Getenv("METRICS_DXTA_DEV_DB_URL"))
	if err != nil {
		return "", err
	}
	defer reposDB.Close()

	dbUrl, err := data.GetReposDbUrl(ctx, reposDB, org, repo)
	if err != nil {
		return "", err
	}

	reposDBCache.Store(cacheKey, dbUrl)

	return dbUrl, nil
}

func NewAPIState(c echo.Context) (APIState, error) {

	ctx := c.Request().Context()

	org := c.Param("org")
	repo := c.Param("repo")

	if org == "" || repo == "" {
		return APIState{}, echo.NewHTTPError(http.StatusBadRequest, "org and repo are required")
	}

	dbUrl, err := getCachedDbUrl(ctx, org, repo)
	if err != nil {
		return APIState{}, err
	}

	driverName := otel.GetDriverName()

	db, err := sql.Open(driverName, dbUrl+"?authToken="+os.Getenv("DXTA_DEV_GROUP_TOKEN"))

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
		org:    org,
		repo:   repo,
		teamId: teamInt,
	}, nil
}

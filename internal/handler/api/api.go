package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/otel"
	"github.com/labstack/echo/v4"

	"github.com/tursodatabase/go-libsql"
)

type APIState struct {
	DB     *sql.DB
	org    string
	repo   string
	teamId *int64
}

var reposDBCache sync.Map

var dbPool sync.Map

func getDirPath() (string, error) {
	dir := filepath.Join(os.TempDir(), "libsql-dir")
	err := os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return "", err
	}

	return dir, nil

}

func getEmbeddedDB(dbUrl, org, repo string) (*libsql.Connector, error) {
	dirPath, err := getDirPath()

	if err != nil {
		return nil, err
	}

	connector, err := libsql.NewEmbeddedReplicaConnector(
		filepath.Join(dirPath, org+"_"+repo),
		dbUrl,
		libsql.WithAuthToken(os.Getenv("DXTA_DEV_GROUP_TOKEN")),
		libsql.WithSyncInterval(time.Minute*5),
	)
	if err != nil {
		return nil, err
	}

	return connector, nil

}

func getDB(ctx context.Context, org, repo string) (*sql.DB, error) {
	cacheKey := org + "/" + repo

	if dbInterface, ok := dbPool.Load(cacheKey); ok {
		return dbInterface.(*sql.DB), nil
	}

	dbUrl, err := getCachedDbUrl(ctx, org, repo)
	if err != nil {
		return nil, err
	}

	connector, err := getEmbeddedDB(dbUrl, org, repo)

	var db *sql.DB

	if err != nil {
		fmt.Println("Using normal db")
		driverName := otel.GetDriverName()

		fullDbUrl := dbUrl + "?authToken=" + os.Getenv("DXTA_DEV_GROUP_TOKEN")

		db, err = sql.Open(driverName, fullDbUrl)
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Using embedded db")
		db = sql.OpenDB(connector)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	dbPool.Store(cacheKey, db)

	return db, nil
}

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

	db, err := getDB(ctx, org, repo)

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

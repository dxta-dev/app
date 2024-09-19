package api

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/dxta-dev/app/internal/data"
	"github.com/labstack/echo/v4"
)

type APIState struct {
	DB     *sql.DB
	org    string
	repo   string
	teamId *int64
}

func NewAPIState(c echo.Context) (APIState, error) {

	org := c.Param("org")
	repo := c.Param("repo")

	if org == "" || repo == "" {
		return APIState{}, echo.NewHTTPError(http.StatusBadRequest, "org and repo are required")
	}

	reposDB, err := sql.Open("libsql", os.Getenv("METRICS_DXTA_DEV_DB_URL"))
	if err != nil {
		return APIState{}, err
	}
	defer reposDB.Close()

	dbUrl, err := data.GetReposDbUrl(reposDB, org, repo)
	if err != nil {
		return APIState{}, err
	}

	db, err := sql.Open("libsql", dbUrl+"?authToken="+os.Getenv("DXTA_DEV_GROUP_TOKEN"))

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
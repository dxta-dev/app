package api

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/data/api"
	"github.com/dxta-dev/app/internal/util"
	"github.com/labstack/echo/v4"
)

func MRSizeHandler(c echo.Context) error {

	org := c.Param("org")
	repo := c.Param("repo")

	reposDB, err := sql.Open("libsql", os.Getenv("METRICS_DXTA_DEV_DB_URL"))
	if err != nil {
		return err
	}
	defer reposDB.Close()

	dbUrl, err := data.GetReposDbUrl(reposDB, org, repo)
	if err != nil {
		return err
	}

	db, err := sql.Open("libsql", dbUrl+"?authToken="+os.Getenv("DXTA_DEV_GROUP_TOKEN"))

	if err != nil {
		return err
	}

	defer db.Close()

	team := c.QueryParam("team")

	var teamInt *int64

	if t, err := strconv.ParseInt(team, 10, 64); err == nil && t > 0 {
		teamInt = &t
	}

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	mrSizes, err := api.GetMRSize(db, context.Background(), org, repo, weeks, teamInt)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mrSizes)
}

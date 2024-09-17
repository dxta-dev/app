package api

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dxta-dev/app/internal/data/api"
	"github.com/dxta-dev/app/internal/util"
	"github.com/labstack/echo/v4"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func MRSizeHandler(c echo.Context) error {

	org := c.Param("org")
	repo := c.Param("repo")
	query := "SELECT db_url FROM repos WHERE organization = ? AND repository = ?"

	reposDB, err := sql.Open("libsql", os.Getenv("METRICS_DXTA_DEV_DB_URL"))

	if err != nil {
		return err
	}

	defer reposDB.Close()

	row := reposDB.QueryRow(query, org, repo)

	var dbUrl string

	err = row.Scan(&dbUrl)

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

	if team == "" {
		teamInt = nil
	} else {
		t, err := strconv.ParseInt(team, 10, 64)
		if err != nil {
			return err
		}
		teamInt = &t
	}

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	mrSize, err := api.GetMRSize(db, context.Background(), org, repo, weeks, teamInt)

	if err != nil {
		return err
	}

	response := fmt.Sprintf("Organization: %s\nRepository: %s\nTeam: %s\nWeeks: %s\n",
		org, repo, team, strings.Join(weeks, ", "))

	for week, value := range mrSize {
		response += fmt.Sprintf("%s, %v\n", week, value)
	}

	fmt.Println(response)

	return nil
}

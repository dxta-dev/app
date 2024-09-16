package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dxta-dev/app/internal/util"
	"github.com/labstack/echo/v4"
)

func CodeChangeHandler(c echo.Context) error {
	org := c.Param("org")
	repo := c.Param("repo")

	team := c.QueryParam("team")

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	response := fmt.Sprintf("Organization: %s\nRepository: %s\nTeam: %s\nWeeks: %s",
		org, repo, team, strings.Join(weeks, ", "))

	return c.String(http.StatusOK, response)
}

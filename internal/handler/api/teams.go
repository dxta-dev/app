package api

import (
	"net/http"

	"github.com/dxta-dev/app/internal/data/api"
	"github.com/labstack/echo/v4"
)

func TeamsHandler(c echo.Context) error {
	ctx := c.Request().Context()

	apiState, err := NewAPIState(c)

	if err != nil {
		return err
	}

	teams, err := api.GetTeams(apiState.DB, ctx)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, teams)
}

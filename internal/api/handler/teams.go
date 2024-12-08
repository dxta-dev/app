package handler

import (
	"net/http"

	"github.com/dxta-dev/app/internal/api"
	"github.com/labstack/echo/v4"
)

func TeamsHandler(c echo.Context) error {
	ctx := c.Request().Context()

	apiState, err := api.NewAPIState(c)

	if err != nil {
		return err
	}

	teams, err := apiState.DB.GetTeams(ctx)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, teams)
}

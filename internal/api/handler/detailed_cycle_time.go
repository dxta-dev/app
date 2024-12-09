package handler

import (
	"net/http"
	"time"

	"github.com/dxta-dev/app/internal/api"
	"github.com/dxta-dev/app/internal/api/data"
	"github.com/dxta-dev/app/internal/util"
	"github.com/labstack/echo/v4"
)

func DetailedCycleTimeHandler(c echo.Context) error {

	ctx := c.Request().Context()

	apiState, err := api.NewAPIState(c)

	if err != nil {
		return err
	}

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	query := data.BuildDetailedCycleTimeQuery(weeks, apiState.TeamId)

	cycleTimes, err := apiState.DB.GetDetailedCycleTime(ctx, query, apiState.Org, apiState.Repo, weeks, apiState.TeamId)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, cycleTimes)
}

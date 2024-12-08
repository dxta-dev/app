package handler

import (
	"net/http"
	"time"

	"github.com/dxta-dev/app/internal/api"
	"github.com/dxta-dev/app/internal/api/data"
	"github.com/dxta-dev/app/internal/util"
	"github.com/labstack/echo/v4"
)

func DeployFrequencyHandler(c echo.Context) error {
	ctx := c.Request().Context()

	apiState, err := api.NewAPIState(c)

	if err != nil {
		return err
	}

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	query := data.BuildDeployFrequencyQuery(weeks)

	result, err := apiState.DB.GetAggregatedValues(ctx, query, apiState.Org, apiState.Repo, weeks, nil)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

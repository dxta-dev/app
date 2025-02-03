package handler

import (
	"net/http"

	"github.com/dxta-dev/app/internal/api"
	"github.com/dxta-dev/app/internal/api/data"
	"github.com/dxta-dev/app/internal/util"
	"github.com/labstack/echo/v4"
)

func CodingTimeHandler(c echo.Context) error {
	ctx := c.Request().Context()

	apiState, err := api.NewAPIState(c)

	if err != nil {
		return err
	}

	startWeek := c.QueryParam("startWeek")
	endWeek := c.QueryParam("endWeek")

	weeks, err := util.GetWeeksRange(startWeek, endWeek)
	if err != nil {
		return err
	}

	query := data.BuildCodingTimeQuery(weeks, apiState.TeamId)

	result, err := apiState.DB.GetAggregatedStatistics(ctx, query, apiState.Org, apiState.Repo, weeks, apiState.TeamId)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

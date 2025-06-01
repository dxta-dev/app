package handler

import (
	"net/http"

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

	weekParam := c.QueryParam("weeks")

	weeksArray := util.GetWeeksArray(weekParam)

	weeksSorted := util.SortISOWeeks(weeksArray)

	query := data.BuildDeployFrequencyQuery(weeksSorted)

	result, err := apiState.DB.GetAggregatedValues(
		ctx,
		query,
		apiState.Org,
		apiState.Repo,
		weeksSorted,
		nil,
	)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

package handler

import (
	"net/http"
	"time"

	"github.com/dxta-dev/app/internal/api/data"
	"github.com/dxta-dev/app/internal/util"
	"github.com/labstack/echo/v4"
)

func DeployTimeHandler(c echo.Context) error {
	ctx := c.Request().Context()

	apiState, err := NewAPIState(c)

	if err != nil {
		return err
	}

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	deployTimes, err := data.GetDeployTime(apiState.DB, ctx, apiState.org, apiState.repo, weeks, apiState.teamId)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, deployTimes)
}

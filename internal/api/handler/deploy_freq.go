package handler

import (
	"net/http"
	"time"

	"github.com/dxta-dev/app/internal/api/data"
	"github.com/dxta-dev/app/internal/util"
	"github.com/labstack/echo/v4"
)

func DeployFrequencyHandler(c echo.Context) error {

	ctx := c.Request().Context()
	apiState, err := NewAPIState(c)

	if err != nil {
		return err
	}

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	deployFrequencies, err := data.GetDeployFrequency(apiState.DB, ctx, apiState.org, apiState.repo, weeks)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, deployFrequencies)
}

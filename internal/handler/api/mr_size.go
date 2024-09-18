package api

import (
	"context"
	"net/http"
	"time"

	"github.com/dxta-dev/app/internal/data/api"
	"github.com/dxta-dev/app/internal/util"
	"github.com/labstack/echo/v4"
)

func MRSizeHandler(c echo.Context) error {

	apiState, err := NewAPIState(c)

	if err != nil {
		return err
	}

	defer apiState.DB.Close()

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	mrSizes, err := api.GetMRSize(apiState.DB, context.Background(), apiState.org, apiState.repo, weeks, apiState.teamId)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mrSizes)
}

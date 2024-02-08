package handlers

import (
	"context"
	"dxta-dev/app/internal/middlewares"
	"dxta-dev/app/internal/templates"
	"strconv"

	"github.com/labstack/echo/v4"

	"dxta-dev/app/internal/data"
)

func (a *App) MergeRequest(c echo.Context) error {
	r := c.Request()
	tenantDatabaseUrl := r.Context().Value(middlewares.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	paramMrId := c.Param("mrid")

	mrId, err := strconv.ParseInt(paramMrId, 10, 64)

	if paramMrId == "" || err != nil {
		return c.String(400, "")
	}

	mrInfo, err := store.GetMergeRequestInfo(mrId)

	if err != nil {
		return err
	}

	events, err := store.GetMergeRequestEvents(mrId)

	if err != nil {
		return err
	}

	components := templates.CircleInfo(events, *mrInfo)

	return components.Render(context.Background(), c.Response().Writer)
}

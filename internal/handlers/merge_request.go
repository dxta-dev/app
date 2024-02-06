package handlers

import (
	"dxta-dev/app/internal/middlewares"
	"fmt"
	"strconv"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"

	"dxta-dev/app/internal/data"
)

func (a *App) MergeRequest(c echo.Context) error {
	r := c.Request()
	tenantDatabaseUrl := r.Context().Value(middlewares.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl:tenantDatabaseUrl,
	}

	paramMrId := c.Param("mrid")

	mrId, err := strconv.ParseInt(paramMrId, 10, 64)

	if paramMrId == "" || err != nil{
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

	var result = struct {
		Mr *data.MergeRequestInfo `json:"mr"`
		Events data.EventSlice `json:"events"`
	}{
		Mr: mrInfo,
		Events: events,
	}

	c.JSON(200, result)
	return nil
}

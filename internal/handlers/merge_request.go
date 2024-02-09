package handlers

import (
	"context"
	"dxta-dev/app/internal/middlewares"
	"dxta-dev/app/internal/templates"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"

	"dxta-dev/app/internal/data"
)

func (a *App) GetMergeRequestInfo(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	tenantDatabaseUrl := r.Context().Value(middlewares.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	paramMrId := c.Param("mrid")
	mrId, err := strconv.ParseInt(paramMrId, 10, 64)

	if paramMrId == "" || err != nil {
		return c.String(400, "")
	}

	parsedURL, err := url.Parse(h.HxCurrentURL)

	var week string
	week = parsedURL.Query().Get("week")

	state := DashboardState{
		week: week,
		mr:   &mrId,
	}

	nextUrl := getNextUrl(state)
	c.Response().Header().Set("HX-Push-Url", nextUrl)

	events, err := store.GetMergeRequestEvents(mrId)

	if err != nil {
		return err
	}

	mergeRequestInfoProps := templates.MergeRequestInfoProps{
		Events: events,
		DeleteEndpoint: fmt.Sprintf("/merge-request/%d", mrId),
		TargetSelector: "#slide-over",
	}

	components := templates.MergeRequestInfo(mergeRequestInfoProps)

	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) RemoveMergeRequestInfo(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	parsedURL, err := url.Parse(h.HxCurrentURL)
	if err != nil {
		return err
	}

	var week string
	week = parsedURL.Query().Get("week")

	state := DashboardState{
		week: week,
		mr:   nil,
	}

	nextUrl := getNextUrl(state)
	c.Response().Header().Set("HX-Push-Url", nextUrl)

	c.NoContent(http.StatusOK)
	return nil
}

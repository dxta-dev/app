package handler

import (
	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/middleware"
	"github.com/dxta-dev/app/internal/otel"
	"github.com/dxta-dev/app/internal/template"

	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) GetMergeRequestInfo(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	ctx := r.Context()
	store := &data.Store{
		DbUrl:      tenantDatabaseUrl,
		DriverName: otel.GetDriverName(),
		Context:    ctx,
	}

	paramMrId := c.Param("mrid")
	mrId, err := strconv.ParseInt(paramMrId, 10, 64)

	if paramMrId == "" || err != nil {
		return c.String(400, "")
	}

	parsedURL, err := url.Parse(h.HxCurrentURL)

	if err != nil {
		return err
	}

	week := parsedURL.Query().Get("week")

	state := DashboardState{
		week: week,
		mr:   &mrId,
	}

	if team := parsedURL.Query().Get("team"); team != "" {
		teamId, err := strconv.ParseInt(team, 10, 64)

		if err != nil {
			return err
		}

		a.State.Team = &teamId
	}

	nextUrl, err := getNextDashboardUrl(a, h.HxCurrentURL, state, nil, true)

	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Push-Url", nextUrl)

	events, uniqueDates, err := store.GetMergeRequestEvents(mrId)

	if err != nil {
		return err
	}

	mergeRequestInfoProps := template.MergeRequestInfoProps{
		Events:         events,
		UniqueDates:    uniqueDates,
		DetailsPageUrl: fmt.Sprintf("/mr/%d", mrId),
		DeleteEndpoint: fmt.Sprintf("/mr-info/%d", mrId),
		TargetSelector: "#mr-info",
	}

	components := template.MergeRequestInfo(mergeRequestInfoProps)

	return components.Render(ctx, c.Response().Writer)
}

func (a *App) RemoveMergeRequestInfo(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	parsedURL, err := url.Parse(h.HxCurrentURL)
	if err != nil {
		return err
	}

	week := parsedURL.Query().Get("week")

	state := DashboardState{
		week: week,
		mr:   nil,
	}

	if team := parsedURL.Query().Get("team"); team != "" {
		teamId, err := strconv.ParseInt(team, 10, 64)

		if err != nil {
			return err
		}

		a.State.Team = &teamId
	}

	nextUrl, err := getNextDashboardUrl(a, h.HxCurrentURL, state, nil, true)

	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Push-Url", nextUrl)

	c.NoContent(http.StatusOK)
	return nil
}

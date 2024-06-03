package handler

import (
	"net/url"
	"strconv"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/middleware"
	"github.com/dxta-dev/app/internal/otel"
	"github.com/dxta-dev/app/internal/template"
	"github.com/dxta-dev/app/internal/util"

	"github.com/labstack/echo/v4"
)

type MergeRequestListRequestState struct {
	team *int64
	week time.Time
}

func loadHxCurrentURLState(h htmx.HxRequestHeader) (*MergeRequestListRequestState, error) {
	var state MergeRequestListRequestState
	state.week = time.Now()
	currentURL, err := url.Parse(h.HxCurrentURL)
	if err != nil {
		return nil, err
	}
	query := currentURL.Query()

	if query.Has("team") {
		value, err := strconv.ParseInt(query.Get("team"), 10, 64)
		if err == nil {
			state.team = &value
		}
	}

	if query.Has("week") {
		date, _, err := util.ParseYearWeek(query.Get("week"))
		if err != nil {
			return nil, err
		}
		state.week = date
	}

	return &state, nil
}

func (a *App) GetMergeRequestWaitingForReviewStack(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	state, err := loadHxCurrentURLState(h)
	if err != nil {
		return err
	}

	a.LoadState(r)

	ctx := r.Context()
	store := &data.Store{
		DbUrl:      tenantDatabaseUrl,
		DriverName: otel.GetDriverName(),
		Context:    ctx,
	}

	var nullRows *data.NullRows

	nullRows, err = store.GetNullRows()

	if err != nil {
		return err
	}

	teamMembers, err := store.GetTeamMembers(state.team)

	if err != nil {
		return err
	}

	var mrStackListProps template.MergeRequestStackedListProps
	mrStackListProps.MergeRequests, err = store.GetMergeRequestWaitingForReviewList(teamMembers, time.Now(), nullRows.UserId)

	if err != nil {
		return err
	}

	components := template.PartialMergeRequestStackedList(mrStackListProps)

	return components.Render(ctx, c.Response().Writer)
}

func (a *App) GetMergeRequestInProgressStack(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	state, err := loadHxCurrentURLState(h)
	if err != nil {
		return err
	}

	a.LoadState(r)

	ctx := r.Context()
	store := &data.Store{
		DbUrl:      tenantDatabaseUrl,
		DriverName: otel.GetDriverName(),
		Context:    ctx,
	}

	var nullRows *data.NullRows

	nullRows, err = store.GetNullRows()

	if err != nil {
		return err
	}

	teamMembers, err := store.GetTeamMembers(state.team)

	if err != nil {
		return err
	}

	var mrStackListProps template.MergeRequestStackedListProps
	mrStackListProps.MergeRequests, err = store.GetMergeRequestInProgressList(time.Now(), teamMembers, nullRows.UserId)

	if err != nil {
		return err
	}

	components := template.PartialMergeRequestStackedList(mrStackListProps)

	return components.Render(ctx, c.Response().Writer)
}

func (a *App) GetMergeRequestReadyToMergeStack(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	state, err := loadHxCurrentURLState(h)
	if err != nil {
		return err
	}

	a.LoadState(r)

	ctx := r.Context()
	store := &data.Store{
		DbUrl:      tenantDatabaseUrl,
		DriverName: otel.GetDriverName(),
		Context:    ctx,
	}

	var nullRows *data.NullRows

	nullRows, err = store.GetNullRows()

	if err != nil {
		return err
	}

	teamMembers, err := store.GetTeamMembers(state.team)

	if err != nil {
		return err
	}

	var mrStackListProps template.MergeRequestStackedListProps
	mrStackListProps.MergeRequests, err = store.GetMergeRequestReadyToMergeList(teamMembers, nullRows.UserId)

	if err != nil {
		return err
	}

	components := template.PartialMergeRequestStackedList(mrStackListProps)

	return components.Render(ctx, c.Response().Writer)
}

func (a *App) GetMergeRequestMergedStack(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	state, err := loadHxCurrentURLState(h)
	if err != nil {
		return err
	}

	a.LoadState(r)

	ctx := r.Context()
	store := &data.Store{
		DbUrl:      tenantDatabaseUrl,
		DriverName: otel.GetDriverName(),
		Context:    ctx,
	}

	var nullRows *data.NullRows

	nullRows, err = store.GetNullRows()

	if err != nil {
		return err
	}

	teamMembers, err := store.GetTeamMembers(state.team)

	if err != nil {
		return err
	}

	var mrStackListProps template.MergeRequestStackedListProps
	mrStackListProps.MergeRequests, err = store.GetMergeRequestMergedList(state.week, teamMembers, nullRows.UserId)

	if err != nil {
		return err
	}

	components := template.PartialMergeRequestStackedList(mrStackListProps)

	return components.Render(ctx, c.Response().Writer)
}

func (a *App) GetMergeRequestClosedStack(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	state, err := loadHxCurrentURLState(h)
	if err != nil {
		return err
	}

	a.LoadState(r)

	ctx := r.Context()
	store := &data.Store{
		DbUrl:      tenantDatabaseUrl,
		DriverName: otel.GetDriverName(),
		Context:    ctx,
	}

	var nullRows *data.NullRows

	nullRows, err = store.GetNullRows()

	if err != nil {
		return err
	}

	teamMembers, err := store.GetTeamMembers(state.team)

	if err != nil {
		return err
	}

	var mrStackListProps template.MergeRequestStackedListProps
	mrStackListProps.MergeRequests, err = store.GetMergeRequestClosedList(state.week, teamMembers, nullRows.UserId)

	if err != nil {
		return err
	}

	components := template.PartialMergeRequestStackedList(mrStackListProps)

	return components.Render(ctx, c.Response().Writer)
}

package handler

import (
	"time"

	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/middleware"
	"github.com/dxta-dev/app/internal/template"

	"context"

	"github.com/labstack/echo/v4"
)

func (a *App) GetMergeRequestInProgressStack(c echo.Context) error {
	r := c.Request()
	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	a.LoadState(r)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	var nullRows *data.NullRows
	var err error

	nullRows, err = store.GetNullRows()

	if err != nil {
		return err
	}

	teamMembers, err := store.GetTeamMembers(a.State.Team)

	if err != nil {
		return err
	}

	date := time.Now() // TODO: week

	var mrStackListProps template.MergeRequestStackedListProps
	mrStackListProps.MergeRequests, err = store.GetMergeRequestInProgressList(date, teamMembers, nullRows.UserId)

	if err != nil {
		return err
	}

	components := template.PartialMergeRequestStackedList(mrStackListProps)

	return components.Render(context.Background(), c.Response().Writer)
}

func (a *App) GetMergeRequestReadyToMergeStack(c echo.Context) error {
	r := c.Request()
	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	a.LoadState(r)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	var nullRows *data.NullRows
	var err error

	nullRows, err = store.GetNullRows()

	if err != nil {
		return err
	}

	teamMembers, err := store.GetTeamMembers(a.State.Team)

	if err != nil {
		return err
	}

	var mrStackListProps template.MergeRequestStackedListProps
	mrStackListProps.MergeRequests, err = store.GetMergeRequestReadyToMergeList(teamMembers, nullRows.UserId)

	if err != nil {
		return err
	}

	components := template.PartialMergeRequestStackedList(mrStackListProps)

	return components.Render(context.Background(), c.Response().Writer)
}

package handlers

import (
	"context"
	"dxta-dev/app/internal/middlewares"
	"dxta-dev/app/internal/templates"
	"dxta-dev/app/internal/utils"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func (a *App) Dashboard(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	tenantDatabaseUrl := r.Context().Value(middlewares.TenantDatabaseURLContext).(string)

	page := &templates.Page{
		Title:   "Dashboard",
		Boosted: h.HxBoosted,
	}

	date, startOfWeek, prevWeek, nextWeek := processWeekPerameters(c, h, tenantDatabaseUrl)

	if h.HxRequest && h.HxTrigger != "" {
		components := templates.SwarmChart(GetSeries(date, tenantDatabaseUrl), startOfWeek)
		return components.Render(context.Background(), c.Response().Writer)
	}

	components := templates.Swarm(page, GetSeries(date, tenantDatabaseUrl), startOfWeek, utils.GetFormattedWeek(date), utils.GetFormattedWeek(time.Now()), prevWeek, nextWeek)

	return components.Render(context.Background(), c.Response().Writer)
}

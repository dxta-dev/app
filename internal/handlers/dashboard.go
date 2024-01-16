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
		Title:   "Charts",
		Boosted: h.HxBoosted,
	}

	date := time.Now()

	weekString := r.URL.Query().Get("week")

	if weekString != "" {
		dateTime, err := utils.ParseYearWeek(weekString)
		if err == nil {
			date = dateTime

			res := c.Response()
			res.Header().Set("HX-Push-Url", "/swarm?week="+weekString)
		}
	}

	startOfWeek := utils.GetStartOfWeek(date)

	if h.HxRequest && h.HxTrigger != "" {
		components := templates.SwarmChart(getSeries(date, tenantDatabaseUrl), startOfWeek)
		return components.Render(context.Background(), c.Response().Writer)
	}

	prevWeek, nextWeek := utils.GetPrevNextWeek(date)

	components := templates.Swarm(page, getSeries(date, tenantDatabaseUrl), startOfWeek, utils.GetFormattedWeek(date), utils.GetFormattedWeek(time.Now()), prevWeek, nextWeek)

	return components.Render(context.Background(), c.Response().Writer)
}

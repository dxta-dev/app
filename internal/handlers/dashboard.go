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
		Requested: h.HxRequest,
	}

	date := time.Now()

	weekString := r.URL.Query().Get("week")

	if weekString != "" {
		dateTime, err := utils.ParseYearWeek(weekString)
		if err == nil {
			date = dateTime

			res := c.Response()
			res.Header().Set("HX-Push-Url", "/dashboard?week="+weekString)
		}
	}

	prevWeek, nextWeek := utils.GetPrevNextWeek(date)

	eventInfo, _ := getData(date, tenantDatabaseUrl)
	weekPickerProps := templates.WeekPickerProps{
		Week: utils.GetFormattedWeek(date),
		CurrentWeek: utils.GetFormattedWeek(time.Now()),
		NextWeek: nextWeek,
		PreviousWeek: prevWeek,
	}

	swarmProps := templates.SwarmProps{
		Series: getSeries(date, tenantDatabaseUrl),
		StartOfTheWeek: utils.GetStartOfTheWeek(date),
	}


	components := templates.DashboardPage(page, swarmProps, weekPickerProps, eventInfo)

	return components.Render(context.Background(), c.Response().Writer)
}

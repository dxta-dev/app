package handlers

import (
	"context"
	"dxta-dev/app/internal/data"
	"dxta-dev/app/internal/graphs"
	"dxta-dev/app/internal/middlewares"
	"dxta-dev/app/internal/templates"
	"dxta-dev/app/internal/utils"
	"sort"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	"github.com/wcharczuk/go-chart/v2/drawing"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func getSwarmSeries(date time.Time, dbUrl string) (templates.SwarmSeries, error) {
	var xvalues []float64
	var yvalues []float64

	store := &data.Store{
		DbUrl: dbUrl,
	}

	events, err := store.GetEventSlices(date)

	if err != nil {
		return templates.SwarmSeries{}, err
	}

	startOfWeek := utils.GetStartOfTheWeek(date)

	var times []time.Time

	sort.Sort(events)

	for _, d := range events {
		t := time.Unix(d.Timestamp/1000, 0)
		times = append(times, t)
	}

	for _, t := range times {
		xSecondsValue := float64(t.Unix() - startOfWeek.Unix())
		xvalues = append(xvalues, xSecondsValue)
		yvalues = append(yvalues, 60*60*12)
	}

	xvalues, yvalues = graphs.Beehive(xvalues, yvalues, 1400, 200, 5)

	colors := []drawing.Color{}

	for i := 0; i < len(xvalues); i++ {
		switch events[i].Type {
		case data.COMMITTED:
			colors = append(colors, drawing.ColorBlue)
		case data.MERGED:
			colors = append(colors, drawing.ColorRed)
		case data.REVIEWED:
			colors = append(colors, drawing.ColorGreen)
		default:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(204, 204, 204, 255))
		}
	}

	return templates.SwarmSeries{
		XValues:   xvalues,
		YValues:   yvalues,
		DotColors: colors,
		Title:     "Swarm",
	}, nil

}

func (a *App) Dashboard(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	tenantDatabaseUrl := r.Context().Value(middlewares.TenantDatabaseURLContext).(string)

	page := &templates.Page{
		Title:     "Charts",
		Boosted:   h.HxBoosted,
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

	weekPickerProps := templates.WeekPickerProps{
		Week:         utils.GetFormattedWeek(date),
		CurrentWeek:  utils.GetFormattedWeek(time.Now()),
		NextWeek:     nextWeek,
		PreviousWeek: prevWeek,
	}

	swarmSeries, err := getSwarmSeries(date, tenantDatabaseUrl)

	if err != nil {
		return err
	}

	swarmProps := templates.SwarmProps{
		Series:         swarmSeries,
		StartOfTheWeek: utils.GetStartOfTheWeek(date),
	}

	components := templates.DashboardPage(page, swarmProps, weekPickerProps)

	return components.Render(context.Background(), c.Response().Writer)
}

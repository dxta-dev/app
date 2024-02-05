package handlers

import (
	"context"
	"dxta-dev/app/internal/data"
	"dxta-dev/app/internal/graphs"
	"dxta-dev/app/internal/middlewares"
	"dxta-dev/app/internal/templates"
	"dxta-dev/app/internal/utils"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	"github.com/wcharczuk/go-chart/v2/drawing"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type DashboardState struct {
	week  string
	event *int64
}

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
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(33, 150, 243, 255)) // Deep Sky Blue
		case data.MERGED:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(156, 39, 176, 255)) // Deep Purple
		case data.CLOSED:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(233, 30, 99, 255)) // Pink
		case data.REVIEWED:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(255, 193, 7, 255)) // Amber
		case data.STARTED_CODING:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(76, 175, 80, 255)) // Green
		case data.ASSIGNED:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(0, 150, 136, 255)) // Teal
		case data.COMMENTED:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(158, 158, 158, 255)) // Grey
		default:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(204, 204, 204, 255)) // Silver
		}
	}

	return templates.SwarmSeries{
		XValues:   xvalues,
		YValues:   yvalues,
		DotColors: colors,
		Title:     "Swarm",
		Events:    events,
	}, nil

}

func getNextUrl(state DashboardState) string {
	params := url.Values{}
	if state.week != "" {
		params.Add("week", state.week)
	}
	if state.event != nil {
		params.Add("event", fmt.Sprintf("%d", *state.event))
	}
	baseUrl := "/dashboard"
	return fmt.Sprintf("%s?%s", baseUrl, params.Encode())
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
	var err error
	var event *int64 = new(int64)
	*event, err = strconv.ParseInt(r.URL.Query().Get("event"), 10, 64)
	if err != nil {
		event = nil
	}
	state := DashboardState{
		week:  r.URL.Query().Get("week"),
		event: event,
	}

	if state.week != "" {
		dateTime, err := utils.ParseYearWeek(state.week)
		if err == nil {
			date = dateTime
		}
	}

	nextUrl := getNextUrl(state)

	c.Response().Header().Set("HX-Push-Url", nextUrl)

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
	var eventIds []int64
	var eventMergeRequestIds []int64
	for _, event := range swarmSeries.Events {
		eventIds = append(eventIds, event.Id)
		eventMergeRequestIds = append(eventMergeRequestIds, event.MergeRequestId)
	}

	swarmProps := templates.SwarmProps{
		Series:               swarmSeries,
		StartOfTheWeek:       utils.GetStartOfTheWeek(date),
		EventIds:             eventIds,
		EventMergeRequestIds: eventMergeRequestIds,
	}

	selectedEvent := data.Event{}
	if event != nil {
		for _, e := range swarmSeries.Events {
			if e.Id == *event {
				selectedEvent = e
			}
		}
	}

	components := templates.DashboardPage(page, swarmProps, weekPickerProps, selectedEvent)

	return components.Render(context.Background(), c.Response().Writer)
}

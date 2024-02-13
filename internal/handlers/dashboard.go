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
	week string
	mr   *int64
}

func absInt64(x int64) int64 {
    if x < 0 {
        return -x
    }
    return x
}

func getSwarmSeries(store *data.Store, date time.Time) (templates.SwarmSeries, error) {
	var xvalues []float64
	var yvalues []float64

	events, err := store.GetEventSlices(date)

	if err != nil {
		return templates.SwarmSeries{}, err
	}

	startOfWeek := utils.GetStartOfTheWeek(date)

	filteredEvents := []data.Event{}

	reviewedMap := make(map[string][]int64)

	sort.Sort(events)

EventLoop:
	for _, e := range events {
		if e.Type == data.REVIEWED {
			key := fmt.Sprintf("%d-%d", e.MergeRequestId, e.Actor)
			reviewed, ok := reviewedMap[key]
			if !ok {
				reviewedMap[key] = []int64{}
			}

			if ok {
				for _, r := range reviewed {
					if absInt64(r-e.Timestamp) < 1000*60*5 {
						continue EventLoop
					}
				}
			}

			reviewedMap[key] = append(reviewedMap[key], e.Timestamp)
		}

		if e.Type == data.COMMITTED ||
			e.Type == data.CLOSED ||
			e.Type == data.REVIEWED ||
			e.Type == data.STARTED_CODING {
			filteredEvents = append(filteredEvents, e)
		}
	}

	events = filteredEvents

	var times []time.Time


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
		//case data.MERGED:
		//	colors = append(colors, drawing.ColorFromAlphaMixedRGBA(156, 39, 176, 255)) // Deep Purple
		case data.CLOSED:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(233, 30, 99, 255)) // Pink
		case data.REVIEWED:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(255, 193, 7, 255)) // Amber
		case data.STARTED_CODING:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(76, 175, 80, 255)) // Green
		//case data.ASSIGNED:
		//	colors = append(colors, drawing.ColorFromAlphaMixedRGBA(0, 150, 136, 255)) // Teal
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
	if state.mr != nil {
		params.Add("mr", fmt.Sprintf("%d", *state.mr))
	}
	baseUrl := "/dashboard"
	encodedParams := params.Encode()
	if encodedParams != "" {
		return fmt.Sprintf("%s?%s", baseUrl, encodedParams)
	}

	return baseUrl
}

func (a *App) Dashboard(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	tenantDatabaseUrl := r.Context().Value(middlewares.TenantDatabaseURLContext).(string)

	page := &templates.Page{
		Title:     "Charts",
		Boosted:   h.HxBoosted,
		Requested: h.HxRequest,
		CacheBust: a.BuildTimestamp,
		DebugMode: a.DebugMode,
	}

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	date := time.Now()
	var err error
	var mr *int64 = new(int64)
	*mr, err = strconv.ParseInt(r.URL.Query().Get("mr"), 10, 64)
	if err != nil {
		mr = nil
	}
	state := DashboardState{
		week: r.URL.Query().Get("week"),
		mr:   mr,
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

	swarmSeries, err := getSwarmSeries(store, date)

	if err != nil {
		return err
	}
	var eventIds []int64
	var eventMergeRequestIds []int64
	var eventOrder []int

	mergeRequestEventOrder := make(map[int64]int)

	for _, event := range swarmSeries.Events {
		eventIds = append(eventIds, event.Id)
		eventMergeRequestIds = append(eventMergeRequestIds, event.MergeRequestId)
		if val, ok := mergeRequestEventOrder[event.MergeRequestId]; ok {
			eventOrder = append(eventOrder, val)
			mergeRequestEventOrder[event.MergeRequestId]++
		} else {
			eventOrder = append(eventOrder, 0)
			mergeRequestEventOrder[event.MergeRequestId] = 1
		}
	}
	swarmProps := templates.SwarmProps{
		Series:               swarmSeries,
		StartOfTheWeek:       utils.GetStartOfTheWeek(date),
		EventIds:             eventIds,
		EventMergeRequestIds: eventMergeRequestIds,
		EventOrder:           eventOrder,
	}

	var mergeRequestInfoProps *templates.MergeRequestInfoProps

	if state.mr != nil {

		events, err := store.GetMergeRequestEvents(*state.mr)

		if err != nil {
			return err
		}

		mergeRequestInfoProps = &templates.MergeRequestInfoProps{
			Events:         events,
			DeleteEndpoint: fmt.Sprintf("/merge-request/%d", *state.mr),
			TargetSelector: "#slide-over",
		}
	}

	components := templates.DashboardPage(page, swarmProps, weekPickerProps, mergeRequestInfoProps)

	return components.Render(context.Background(), c.Response().Writer)
}

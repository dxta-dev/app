package handler

import (
	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/graph"
	"github.com/dxta-dev/app/internal/middleware"
	"github.com/dxta-dev/app/internal/template"
	"github.com/dxta-dev/app/internal/util"

	"context"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	_ "modernc.org/sqlite"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	_ "github.com/libsql/libsql-client-go/libsql"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type DashboardState struct {
	week string
	mr   *int64
	team *int64
}

func getSwarmSeries(store *data.Store, date time.Time, teamMembers data.TeamMembers) (template.SwarmSeries, error) {
	var xvalues []float64
	var yvalues []float64

	events, err := store.GetEventSlices(date, teamMembers)

	if err != nil {
		return template.SwarmSeries{}, err
	}

	startOfWeek := util.GetStartOfTheWeek(date)

	filteredEvents := []data.Event{}

	for _, e := range events {
		if e.Type == data.COMMITTED ||
			e.Type == data.CLOSED ||
			e.Type == data.REVIEWED ||
			e.Type == data.STARTED_CODING {
			filteredEvents = append(filteredEvents, e)
		}
	}

	events = filteredEvents

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

	xvalues, yvalues = graph.Beehive(xvalues, yvalues, 1400, 200, 5)

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

	return template.SwarmSeries{
		XValues:   xvalues,
		YValues:   yvalues,
		DotColors: colors,
		Title:     "Swarm",
		Events:    events,
	}, nil

}

func getNextDashboardUrl(currentUrl string, state DashboardState) (string, error) {
	params := url.Values{}

	parsedURL, err := url.Parse(currentUrl)

	if err != nil {
		return "", err
	}

	requestUri := parsedURL.Path

	if state.week != "" {
		params.Add("week", state.week)
	}
	if state.mr != nil {
		params.Add("mr", fmt.Sprintf("%d", *state.mr))
	}
	if state.team != nil {
		params.Add("team", fmt.Sprint(*state.team))
	}
	encodedParams := params.Encode()
	if encodedParams != "" {
		return fmt.Sprintf("%s?%s", requestUri, encodedParams), nil
	}

	return requestUri, nil
}

func (a *App) DashboardPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	page := &template.Page{
		Title:     "Dashboard - DXTA",
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

	var team *int64
	if r.URL.Query().Has("team") {
		value, err := strconv.ParseInt(r.URL.Query().Get("team"), 10, 64)
		if err == nil {
			team = &value
		}
	}

	teamMembers, err := store.GetTeamMembers(team)

	if err != nil {
		return err
	}

	teams, err := store.GetTeams()

	if err != nil {
		return err
	}

	state := DashboardState{
		week: r.URL.Query().Get("week"),
		mr:   mr,
		team: team,
	}

	if state.week != "" {
		dateTime, err := util.ParseYearWeek(state.week)
		if err == nil {
			date = dateTime
		}
	}

	var nextUrl string

	if h.HxRequest && !h.HxBoosted {
		nextUrl, err = getNextDashboardUrl(h.HxCurrentURL, state)
		if err != nil {
			return err
		}
	} else {
		nextUrl, err = getNextDashboardUrl(r.URL.RequestURI(), state)
		if err != nil {
			return err
		}
	}

	c.Response().Header().Set("HX-Push-Url", nextUrl)

	searchParams := url.Values{}
	if team != nil {
		searchParams.Set("team", fmt.Sprint(*team))
	}
	if state.week != "" {
		searchParams.Set("week", state.week)
	}

	prevWeek, nextWeek := util.GetPrevNextWeek(date)

	weekPickerProps := template.WeekPickerProps{
		Week:         util.GetFormattedWeek(date),
		SearchParams: searchParams,
		PreviousWeek: prevWeek,
		CurrentWeek:  util.GetFormattedWeek(time.Now()),
		NextWeek:     nextWeek,
	}

	teamPickerProps := template.TeamPickerProps{
		Teams:        teams,
		SelectedTeam: team,
		SearchParams: searchParams,
		BaseUrl:      "/",
	}

	swarmSeries, err := getSwarmSeries(store, date, teamMembers)

	if err != nil {
		return err
	}
	var eventIds []int64
	var eventMergeRequestIds []int64
	for _, event := range swarmSeries.Events {
		eventIds = append(eventIds, event.Id)
		eventMergeRequestIds = append(eventMergeRequestIds, event.MergeRequestId)
	}

	swarmProps := template.SwarmProps{
		Series:               swarmSeries,
		StartOfTheWeek:       util.GetStartOfTheWeek(date),
		EventIds:             eventIds,
		EventMergeRequestIds: eventMergeRequestIds,
	}

	var mergeRequestInfoProps *template.MergeRequestInfoProps

	if state.mr != nil {

		events, err := store.GetMergeRequestEvents(*state.mr)

		if err != nil {
			return err
		}

		mergeRequestInfoProps = &template.MergeRequestInfoProps{
			Events:         events,
			DeleteEndpoint: fmt.Sprintf("/merge-request/%d", *state.mr),
			TargetSelector: "#slide-over",
		}
	}

	components := template.DashboardPage(page, swarmProps, weekPickerProps, mergeRequestInfoProps, teamPickerProps)

	return components.Render(context.Background(), c.Response().Writer)
}

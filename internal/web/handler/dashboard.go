package handler

import (
	"sync"

	"github.com/dxta-dev/app/internal/util"
	"github.com/dxta-dev/app/internal/web/data"
	"github.com/dxta-dev/app/internal/web/graph"
	"github.com/dxta-dev/app/internal/web/middleware"
	"github.com/dxta-dev/app/internal/web/template"

	"fmt"
	"net/url"
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
}

func getSwarmSeries(store *data.Store, date time.Time, teamMembers []int64) (template.SwarmSeries, error) {
	var xvalues []float64
	var yvalues []float64

	events, err := store.GetEventSlices(date, teamMembers)

	if err != nil {
		return template.SwarmSeries{}, err
	}

	startOfWeek := util.GetStartOfTheWeek(date)

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

	xvalues, yvalues = graph.Beehive(xvalues, yvalues, 1400, 200, 5)

	colors := []drawing.Color{}

	for i := 0; i < len(xvalues); i++ {
		switch events[i].Type {
		case data.COMMITTED:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(33, 150, 243, 255)) // Deep Sky Blue
		case data.CLOSED:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(233, 30, 99, 255)) // Pink
		case data.REVIEWED:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(255, 193, 7, 255)) // Amber
		case data.STARTED_CODING:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(76, 175, 80, 255)) // Green
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

func getNextDashboardUrl(app *App, currentUrl string, state DashboardState, params url.Values, includeAppState bool) (string, error) {
	if params == nil {
		params = url.Values{}
	}

	parsedURL, err := url.Parse(currentUrl)

	if err != nil {
		return "", err
	}

	requestUri := parsedURL.Path

	if state.week != "" && !params.Has("week") {
		params.Add("week", state.week)
	}

	if state.mr != nil && !params.Has("mr") {
		params.Add("mr", fmt.Sprintf("%d", *state.mr))
	}

	var nextUrl string
	if includeAppState {
		nextUrl, err = app.GetUrlAppState(requestUri, params)
	} else {
		nextUrl, err = GetUrl(requestUri, params)
	}

	if err != nil {
		return "", err
	}

	return nextUrl, nil
}

func (a *App) DashboardPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	var err error

	ctx := c.Request().Context()
	store := r.Context().Value(middleware.StoreContextKey).(*data.Store)

	a.GenerateNonce()
	a.LoadState(r)

	timeNow := time.Now()
	date := time.Now()
	var mr *int64 = new(int64)
	*mr, err = strconv.ParseInt(r.URL.Query().Get("mr"), 10, 64)
	if err != nil {
		mr = nil
	}

	team := a.State.Team

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
	}

	if state.week != "" {
		dateTime, _, err := util.ParseYearWeek(state.week)
		if err == nil {
			date = dateTime
		}
	}

	var nextUrl string

	if h.HxRequest && !h.HxBoosted {
		nextUrl, err = getNextDashboardUrl(a, r.URL.Path, state, nil, true)
		if err != nil {
			return err
		}
	} else {
		nextUrl, err = getNextDashboardUrl(a, r.URL.Path, state, nil, true)
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
	queriedWeek := util.GetFormattedWeek(date)
	currentWeek := util.GetFormattedWeek(time.Now())

	prevWeekParams := url.Values{}
	prevWeekParams.Set("week", prevWeek)
	previousWeekUrl, err := getNextDashboardUrl(a, r.URL.Path, DashboardState{mr: nil, week: state.week}, prevWeekParams, true)

	if err != nil {
		return err
	}

	nextWeekParams := url.Values{}
	nextWeekParams.Set("week", nextWeek)
	nextWeekUrl, err := getNextDashboardUrl(a, r.URL.Path, DashboardState{mr: nil, week: state.week}, nextWeekParams, true)

	if err != nil {
		return err
	}

	if queriedWeek == currentWeek {
		nextWeekUrl = ""
	}

	currentWeekParams := url.Values{}
	currentWeekParams.Set("week", currentWeek)
	currentWeekUrl, err := getNextDashboardUrl(a, r.URL.Path, DashboardState{mr: nil, week: state.week}, currentWeekParams, true)

	if err != nil {
		return err
	}

	startEndWeekDays, err := util.GetStartEndWeekDates(queriedWeek)

	if err != nil {
		return err
	}

	weekPickerProps := template.WeekPickerProps{
		Week:              queriedWeek,
		StartEndWeekDates: startEndWeekDays,
		CurrentWeek:       util.GetFormattedWeek(time.Now()),
		PreviousWeekUrl:   previousWeekUrl,
		CurrentWeekUrl:    currentWeekUrl,
		NextWeekUrl:       nextWeekUrl,
	}

	var templTeams []template.Team

	for _, team := range teams {
		params := url.Values{}
		params.Set("team", fmt.Sprint(team.Id))
		teamUrl, err := getNextDashboardUrl(a, r.URL.Path, state, params, true)
		if err != nil {
			return err
		}
		templTeams = append(templTeams, template.Team{
			Id:   team.Id,
			Name: team.Name,
			Url:  teamUrl,
		})
	}

	noTeamUrl, err := getNextDashboardUrl(a, r.URL.Path, DashboardState{week: state.week, mr: state.mr}, nil, false)

	if err != nil {
		return err
	}

	teamPickerProps := template.TeamPickerProps{
		Teams:        templTeams,
		SelectedTeam: team,
		NoTeamUrl:    noTeamUrl,
	}

	swarmSeries, err := getSwarmSeries(store, date, teamMembers)

	if err != nil {
		return err
	}
	var eventIds []int64
	var eventMergeRequestIds []int64
	var mergeRequestsStartedThisWeek = make(map[int64]bool)
	var mergeRequestsClosedThisWeek = make(map[int64]bool)
	for _, event := range swarmSeries.Events {
		eventIds = append(eventIds, event.Id)
		if event.Type == data.STARTED_CODING {
			mergeRequestsStartedThisWeek[event.MergeRequestId] = true
		}
		if event.Type == data.MERGED || event.Type == data.CLOSED || queriedWeek == currentWeek {
			mergeRequestsClosedThisWeek[event.MergeRequestId] = true
		}
		eventMergeRequestIds = append(eventMergeRequestIds, event.MergeRequestId)
	}
	startedMergeRequestIds := []int64{}
	closedMergeRequestIds := []int64{}

	for id, _ := range mergeRequestsStartedThisWeek {
		startedMergeRequestIds = append(startedMergeRequestIds, id)
	}
	for id, _ := range mergeRequestsClosedThisWeek {
		closedMergeRequestIds = append(closedMergeRequestIds, id)
	}

	swarmProps := template.SwarmProps{
		Series:                 swarmSeries,
		StartOfTheWeek:         util.GetStartOfTheWeek(date),
		EventIds:               eventIds,
		EventMergeRequestIds:   eventMergeRequestIds,
		StartedMergeRequestIds: startedMergeRequestIds,
		ClosedMergeRequestIds:  closedMergeRequestIds,
	}

	var mergeRequestInfoProps *template.MergeRequestInfoProps

	if state.mr != nil {

		events, uniqueDates, err := store.GetMergeRequestEvents(*state.mr)

		if err != nil {
			return err
		}

		mergeRequestInfoProps = &template.MergeRequestInfoProps{
			Events:           events,
			UniqueDates:      uniqueDates,
			DetailsPageUrl:   fmt.Sprintf("/mr/%d", *state.mr),
			ShouldOpenMrInfo: !(h.HxBoosted && h.HxRequest) && state.mr != nil,
			TargetSelector:   "#mr-info",
		}
	}

	navState, err := a.GetNavState()

	if err != nil {
		return err
	}

	var nullRows *data.NullRows

	nullRows, err = store.GetNullRows()

	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	errCh := make(chan error, 5)

	var (
		mergeRequestsWaitingForReview template.MergeRequestStackedListProps
		mergeRequestsReadyToMerge     template.MergeRequestStackedListProps
		mergeRequestsInProgress       template.MergeRequestStackedListProps
		mergeRequestsMerged           template.MergeRequestStackedListProps
		mergeRequestsClosed           template.MergeRequestStackedListProps
		mergeRequestsStale            template.MergeRequestStackedListProps
		mrsInReview                   []data.MergeRequestListItemData
	)

	wg.Add(2)

	isQueryCurrentWeek := queriedWeek == util.GetFormattedWeek(timeNow)
	if isQueryCurrentWeek {
		wg.Add(3)
		go func() {
			defer wg.Done()
			mergeRequestsInProgress.MergeRequests, err = store.GetMergeRequestInProgressList(timeNow, teamMembers, nullRows.UserId)
			if err != nil {
				errCh <- err
				return
			}
		}()

		mergeRequestsInProgress.Id = "mrs-in-progress"
		mergeRequestsInProgress.Title = "In progress"
		mergeRequestsInProgress.MRStatusIconProps = template.MRInProgressIconProps

		go func() {
			defer wg.Done()
			mergeRequestsReadyToMerge.MergeRequests, err = store.GetMergeRequestReadyToMergeList(teamMembers, nullRows.UserId, queriedWeek)
			if err != nil {
				errCh <- err
				return
			}
		}()

		mergeRequestsReadyToMerge.Id = "mrs-ready-to-merge"
		mergeRequestsReadyToMerge.Title = "Ready to merge"
		mergeRequestsReadyToMerge.MRStatusIconProps = template.MRReadyToBeMergedIconProps

		go func() {
			defer wg.Done()
			mrsInReview, err = store.GetMergeRequestWaitingForReviewList(teamMembers, timeNow, nullRows.UserId)
			if err != nil {
				errCh <- err
				return
			}
		}()

		var recentMrsWaitingForReview []data.MergeRequestListItemData
		var staleMrsWaitingForReview []data.MergeRequestListItemData
		var lastMonday = util.GetStartOfTheWeek(timeNow).Add(-1 * 7 * 24 * time.Hour)
		for _, mr := range mrsInReview {
			if mr.LastEventAt.Before(lastMonday) {
				staleMrsWaitingForReview = append(staleMrsWaitingForReview, mr)
			} else {
				recentMrsWaitingForReview = append(recentMrsWaitingForReview, mr)
			}
		}

		mergeRequestsWaitingForReview.Id = "mrs-waiting-review"
		mergeRequestsWaitingForReview.Title = "Waiting for review"
		mergeRequestsWaitingForReview.MergeRequests = recentMrsWaitingForReview
		mergeRequestsWaitingForReview.MRStatusIconProps = template.MRWaitingToBeReviewedIconProps

		mergeRequestsStale.Id = "mrs-stale"
		mergeRequestsStale.Title = "Stale"
		mergeRequestsStale.MergeRequests = staleMrsWaitingForReview
		mergeRequestsStale.MRStatusIconProps = template.MRStaleIconProps

	}

	go func() {
		defer wg.Done()
		mergeRequestsMerged.MergeRequests, err = store.GetMergeRequestMergedList(date, teamMembers, nullRows.UserId)
		if err != nil {
			errCh <- err
			return
		}
	}()

	mergeRequestsMerged.Id = "mrs-merged"
	mergeRequestsMerged.Title = "Merged"
	mergeRequestsMerged.MRStatusIconProps = template.MRMergedIconProps

	go func() {
		defer wg.Done()
		mergeRequestsClosed.MergeRequests, err = store.GetMergeRequestClosedList(date, teamMembers, nullRows.UserId)
		if err != nil {
			errCh <- err
			return
		}
	}()

	mergeRequestsClosed.Id = "mrs-closed"
	mergeRequestsClosed.Title = "Closed"
	mergeRequestsClosed.MRStatusIconProps = template.MRClosedIconProps

	wg.Wait()

	close(errCh)
	for e := range errCh {
		if e != nil {
			err = e
			break
		}
	}

	page := &template.Page{
		RouteId:   "/",
		Title:     "Dashboard - DXTA",
		Boosted:   h.HxBoosted,
		Requested: h.HxRequest,
		CacheBust: a.BuildTimestamp,
		DebugMode: a.DebugMode,
		NavState:  navState,
		Nonce:     a.Nonce,
	}

	components := template.DashboardPage(page,
		swarmProps,
		weekPickerProps,
		mergeRequestInfoProps,
		teamPickerProps,
		isQueryCurrentWeek,
		mergeRequestsClosed,
		mergeRequestsMerged,
		mergeRequestsInProgress,
		mergeRequestsReadyToMerge,
		mergeRequestsWaitingForReview,
		mergeRequestsStale)

	return components.Render(ctx, c.Response().Writer)
}

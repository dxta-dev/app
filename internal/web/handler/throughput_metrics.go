package handler

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/dxta-dev/app/internal/util"
	"github.com/dxta-dev/app/internal/web/data"
	"github.com/dxta-dev/app/internal/web/middleware"
	"github.com/dxta-dev/app/internal/web/template"

	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) ThroughputMetricsPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	store := r.Context().Value(middleware.StoreContextKey).(*data.Store)

	a.GenerateNonce()
	a.LoadState(r)

	ctx := r.Context()

	teams, err := store.GetTeams()

	if err != nil {
		return err
	}

	team := a.State.Team

	teamMembers, err := store.GetTeamMembers(team)

	if err != nil {
		return err
	}

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	var wg sync.WaitGroup

	errCh := make(chan error, 6)

	var (
		totalCommits                    map[string]data.CommitCountByWeek
		averageTotalCommitsByNWeeks     float64
		totalMrsOpened                  map[string]data.MrCountByWeek
		averageMrsOpenedByNWeeks        float64
		mergeFrequency                  map[string]data.MergeFrequencyByWeek
		averageMergeFrequencyByNWeeks   float64
		totalReviews                    map[string]data.TotalReviewsByWeek
		averageReviewsByNWeeks          float64
		totalCodeChanges                map[string]data.CodeChangesCount
		averageTotalCodeChangesByNWeeks float64
		deployFrequency                 map[string]data.DeployFrequencyByWeek
		averageDeployFrequencyByNWeeks  float64
	)

	wg.Add(6)

	go func() {
		defer wg.Done()
		totalCommits, averageTotalCommitsByNWeeks, err = store.GetTotalCommits(weeks, teamMembers)
		if err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		totalMrsOpened, averageMrsOpenedByNWeeks, err = store.GetTotalMrsOpened(weeks, teamMembers)
		if err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		mergeFrequency, averageMergeFrequencyByNWeeks, err = store.GetMergeFrequency(weeks, teamMembers)
		if err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		totalReviews, averageReviewsByNWeeks, err = store.GetTotalReviews(weeks, teamMembers)
		if err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		totalCodeChanges, averageTotalCodeChangesByNWeeks, err = store.GetTotalCodeChanges(weeks, teamMembers)
		if err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		deployFrequency, averageDeployFrequencyByNWeeks, err = store.GetDeployFrequency(weeks)
		if err != nil {
			errCh <- err
			return
		}
	}()

	wg.Wait()

	close(errCh)
	for e := range errCh {
		if e != nil {
			err = e
			break
		}
	}

	if err != nil {
		return err
	}

	totalCommitsXValues := make([]float64, len(weeks))
	totalCommitsYValues := make([]float64, len(weeks))
	startEndWeek := make([]template.StartEndWeek, len(weeks))

	for i, week := range weeks {
		totalCommitsXValues[i] = float64(i)
		totalCommitsYValues[i] = float64(totalCommits[week].Count)
		startWeek, endWeek, err := util.ParseYearWeek(week)
		if err != nil {
			return err
		}
		startEndWeek[i] = template.StartEndWeek{
			Start: startWeek.Format("Jan 02"),
			End:   endWeek.Format("Jan 02"),
		}
	}

	formattedTotalCommitsYValues := make([]string, len(totalCommitsYValues))

	for i, value := range totalCommitsYValues {
		formattedTotalCommitsYValues[i] = util.FormatYAxisValues(value)
	}

	averageTotalCommitsSeries := template.TimeSeries{
		Title:   "Total Commits",
		XValues: totalCommitsXValues,
		YValues: totalCommitsYValues,
		Weeks:   weeks,
	}

	averageTotalCommitsSeriesProps := template.TimeSeriesProps{
		Series:           averageTotalCommitsSeries,
		StartEndWeeks:    startEndWeek,
		FormattedYValues: formattedTotalCommitsYValues,
		InfoText:         fmt.Sprintf("AVG Commits per week: %v", util.FormatYAxisValues(averageTotalCommitsByNWeeks)),
	}

	if err != nil {
		return err
	}

	totalMrsOpenedXValues := make([]float64, len(weeks))
	totalMrsOpenedYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		totalMrsOpenedXValues[i] = float64(i)
		totalMrsOpenedYValues[i] = float64(totalMrsOpened[week].Count)
	}

	formattedTotalMrsOpenedYValues := make([]string, len(totalMrsOpenedYValues))

	for i, value := range totalMrsOpenedYValues {
		formattedTotalMrsOpenedYValues[i] = util.FormatYAxisValues(value)
	}

	averageMrsOpenedSeries := template.TimeSeries{
		Title:   "Total MRs Opened",
		XValues: totalMrsOpenedXValues,
		YValues: totalMrsOpenedYValues,
		Weeks:   weeks,
	}

	averageMrsOpenedSeriesProps := template.TimeSeriesProps{
		Series:           averageMrsOpenedSeries,
		StartEndWeeks:    startEndWeek,
		FormattedYValues: formattedTotalMrsOpenedYValues,
		InfoText:         fmt.Sprintf("AVG MRs Opened per week: %v", util.FormatYAxisValues(averageMrsOpenedByNWeeks)),
	}

	if err != nil {
		return err
	}

	mergeFrequencyXValues := make([]float64, len(weeks))
	mergeFrequencyYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		mergeFrequencyXValues[i] = float64(i)
		mergeFrequencyYValues[i] = float64(mergeFrequency[week].Amount)
	}

	formattedMergeFrequencyYValues := make([]string, len(mergeFrequencyYValues))

	for i, value := range mergeFrequencyYValues {
		formattedMergeFrequencyYValues[i] = util.FormatYAxisValues(value)
	}

	averageMergeFrequencySeries := template.TimeSeries{
		Title:   "Merge Frequency",
		XValues: mergeFrequencyXValues,
		YValues: mergeFrequencyYValues,
		Weeks:   weeks,
	}

	averageMergeFrequencySeriesProps := template.TimeSeriesProps{
		Series:           averageMergeFrequencySeries,
		StartEndWeeks:    startEndWeek,
		FormattedYValues: formattedMergeFrequencyYValues,
		InfoText:         fmt.Sprintf("AVG Merge Frequency per week: %v", util.FormatYAxisValues(averageMergeFrequencyByNWeeks)),
	}

	if err != nil {
		return err
	}

	totalReviewsXValues := make([]float64, len(weeks))
	totalReviewsYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		totalReviewsXValues[i] = float64(i)
		totalReviewsYValues[i] = float64(totalReviews[week].Count)
	}

	formattedTotalReviewsYValues := make([]string, len(totalReviewsYValues))

	for i, value := range totalReviewsYValues {
		formattedTotalReviewsYValues[i] = util.FormatYAxisValues(value)
	}

	averageReviewsSeries := template.TimeSeries{
		Title:   "Total Reviews",
		XValues: totalReviewsXValues,
		YValues: totalReviewsYValues,
		Weeks:   weeks,
	}

	averageReviewsSeriesProps := template.TimeSeriesProps{
		Series:           averageReviewsSeries,
		StartEndWeeks:    startEndWeek,
		FormattedYValues: formattedTotalReviewsYValues,
		InfoText:         fmt.Sprintf("AVG Total Reviews per week: %v", util.FormatYAxisValues(averageReviewsByNWeeks)),
	}

	if err != nil {
		return err
	}

	totalCodeChangesXValues := make([]float64, len(weeks))
	totalCodeChangesYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		totalCodeChangesXValues[i] = float64(i)
		totalCodeChangesYValues[i] = float64(totalCodeChanges[week].Count)
	}

	formattedTotalCodeChangesYValues := make([]string, len(totalCodeChangesYValues))

	for i, value := range totalCodeChangesYValues {
		formattedTotalCodeChangesYValues[i] = util.FormatYAxisValues(value)
	}

	averageCodeChangesSeries := template.TimeSeries{
		Title:   "Total Code Changes",
		XValues: totalCodeChangesXValues,
		YValues: totalCodeChangesYValues,
		Weeks:   weeks}

	averageTotalCodeChangesProps := template.TimeSeriesProps{
		Series:           averageCodeChangesSeries,
		StartEndWeeks:    startEndWeek,
		FormattedYValues: formattedTotalCodeChangesYValues,
		InfoText:         fmt.Sprintf("AVG Total Code Changes per week: %v", util.FormatYAxisValues(averageTotalCodeChangesByNWeeks)),
	}

	if err != nil {
		return err
	}

	deployFrequencyXValues := make([]float64, len(weeks))
	deployFrequencyYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		deployFrequencyXValues[i] = float64(i)
		deployFrequencyYValues[i] = float64(deployFrequency[week].Amount)
	}

	formattedDeployFrequencyYValues := make([]string, len(deployFrequencyYValues))

	for i, value := range deployFrequencyYValues {
		formattedDeployFrequencyYValues[i] = util.FormatYAxisValues(value)
	}

	averageDeployFrequencySeries := template.TimeSeries{
		Title:   "Deploy Frequency",
		XValues: deployFrequencyXValues,
		YValues: deployFrequencyYValues,
		Weeks:   weeks,
	}

	averageDeployFrequencySeriesProps := template.TimeSeriesProps{
		Series:           averageDeployFrequencySeries,
		StartEndWeeks:    startEndWeek,
		FormattedYValues: formattedDeployFrequencyYValues,
		InfoText:         fmt.Sprintf("AVG Deploy Frequency per week: %v", util.FormatYAxisValues(averageDeployFrequencyByNWeeks)),
	}

	props := template.ThroughputMetricsProps{
		TotalCommitsSeriesProps:     averageTotalCommitsSeriesProps,
		TotalMrsOpenedSeriesProps:   averageMrsOpenedSeriesProps,
		MergeFrequencySeriesProps:   averageMergeFrequencySeriesProps,
		TotalReviewsSeriesProps:     averageReviewsSeriesProps,
		TotalCodeChangesSeriesProps: averageTotalCodeChangesProps,
		DeployFrequencySeriesProps:  averageDeployFrequencySeriesProps,
	}

	var templTeams []template.Team

	for _, team := range teams {
		params := url.Values{}
		params.Set("team", fmt.Sprint(team.Id))
		teamUrl, err := a.GetUrlAppState(r.URL.Path, params)
		if err != nil {
			return err
		}
		templTeams = append(templTeams, template.Team{
			Id:   team.Id,
			Name: team.Name,
			Url:  teamUrl,
		})
	}

	teamPickerProps := template.TeamPickerProps{
		Teams:        templTeams,
		SelectedTeam: team,
		NoTeamUrl:    r.URL.Path,
	}

	navState, err := a.GetNavState()

	if err != nil {
		return err
	}

	page := &template.Page{
		RouteId:   "/metrics/throughput",
		Title:     "Throughput Metrics - DXTA",
		Boosted:   h.HxBoosted,
		CacheBust: a.BuildTimestamp,
		DebugMode: a.DebugMode,
		NavState:  navState,
		Nonce:     a.Nonce,
	}

	components := template.ThroughputMetricsPage(page, props, teamPickerProps)
	return components.Render(ctx, c.Response().Writer)
}

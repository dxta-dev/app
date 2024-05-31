package handler

import (
	"fmt"
	"net/url"

	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/middleware"
	"github.com/dxta-dev/app/internal/otel"
	"github.com/dxta-dev/app/internal/template"
	"github.com/dxta-dev/app/internal/util"

	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) ThroughputMetricsPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	a.GenerateNonce()
	a.LoadState(r)

	currentTime := time.Now()

	numWeeksAgo := currentTime.Add(-12 * 7 * 24 * time.Hour)

	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	ctx := r.Context()
	store := &data.Store{
		DbUrl:      tenantDatabaseUrl,
		DriverName: otel.GetDriverName(),
		Context:    ctx,
	}

	crawlInstances, err := store.GetCrawlInstances(numWeeksAgo.Unix(), currentTime.Unix())
	if err != nil {
		return err
	}

	cutOffWeeks := data.GetCutOffWeeks(crawlInstances)

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

	totalCommits, averageTotalCommitsByNWeeks, err := store.GetTotalCommits(weeks, teamMembers)

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
		CutOffWeeks:      cutOffWeeks,
		FormattedYValues: formattedTotalCommitsYValues,
		InfoText:         fmt.Sprintf("AVG Commits per week: %v", util.FormatYAxisValues(averageTotalCommitsByNWeeks)),
	}

	totalMrsOpened, averageMrsOpenedByNWeeks, err := store.GetTotalMrsOpened(weeks, teamMembers)

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
		CutOffWeeks:      cutOffWeeks,
		FormattedYValues: formattedTotalMrsOpenedYValues,
		InfoText:         fmt.Sprintf("AVG MRs Opened per week: %v", util.FormatYAxisValues(averageMrsOpenedByNWeeks)),
	}

	mergeFrequency, averageMergeFrequencyByNWeeks, err := store.GetMergeFrequency(weeks, teamMembers)

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
		CutOffWeeks:      cutOffWeeks,
		FormattedYValues: formattedMergeFrequencyYValues,
		InfoText:         fmt.Sprintf("AVG Merge Frequency per week: %v", util.FormatYAxisValues(averageMergeFrequencyByNWeeks)),
	}

	totalReviews, averageReviewsByNWeeks, err := store.GetTotalReviews(weeks, teamMembers)

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
		CutOffWeeks:      cutOffWeeks,
		FormattedYValues: formattedTotalReviewsYValues,
		InfoText:         fmt.Sprintf("AVG Total Reviews per week: %v", util.FormatYAxisValues(averageReviewsByNWeeks)),
	}

	totalCodeChanges, averageTotalCodeChangesByNWeeks, err := store.GetTotalCodeChanges(weeks, teamMembers)

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
		CutOffWeeks:      cutOffWeeks,
		FormattedYValues: formattedTotalCodeChangesYValues,
		InfoText:         fmt.Sprintf("AVG Total Code Changes per week: %v", util.FormatYAxisValues(averageTotalCodeChangesByNWeeks)),
	}
	if err != nil {
		return err
	}

	props := template.ThroughputMetricsProps{
		TotalCommitsSeriesProps:     averageTotalCommitsSeriesProps,
		TotalMrsOpenedSeriesProps:   averageMrsOpenedSeriesProps,
		MergeFrequencySeriesProps:   averageMergeFrequencySeriesProps,
		TotalReviewsSeriesProps:     averageReviewsSeriesProps,
		TotalCodeChangesSeriesProps: averageTotalCodeChangesProps,
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

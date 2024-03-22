package handler

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/middleware"
	"github.com/dxta-dev/app/internal/template"
	"github.com/dxta-dev/app/internal/util"

	"context"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) ThroughputMetricsPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &template.Page{
		Title:     "Throughput Metrics - DXTA",
		Boosted:   h.HxBoosted,
		CacheBust: a.BuildTimestamp,
		DebugMode: a.DebugMode,
	}

	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	teams, err := store.GetTeams()

	if err != nil {
		return err
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
			Start: startWeek,
			End:   endWeek,
		}
	}

	averageTotalCommitsSeries := template.TimeSeries{
		Title:   "Total Commits",
		XValues: totalCommitsXValues,
		YValues: totalCommitsYValues,
		Weeks:   weeks,
	}

	averageTotalCommitsSeriesProps := template.TimeSeriesProps{
		Series:        averageTotalCommitsSeries,
		StartEndWeeks: startEndWeek,
		InfoText:      fmt.Sprintf("AVG Commits per week: %v", util.FormatYAxisValues(averageTotalCommitsByNWeeks)),
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

	averageMrsOpenedSeries := template.TimeSeries{
		Title:   "Total MRs Opened",
		XValues: totalMrsOpenedXValues,
		YValues: totalMrsOpenedYValues,
		Weeks:   weeks,
	}

	averageMrsOpenedSeriesProps := template.TimeSeriesProps{
		Series:        averageMrsOpenedSeries,
		StartEndWeeks: startEndWeek,
		InfoText:      fmt.Sprintf("AVG MRs Opened per week: %v", util.FormatYAxisValues(averageMrsOpenedByNWeeks)),
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

	averageMergeFrequencySeries := template.TimeSeries{
		Title:   "Merge Frequency",
		XValues: mergeFrequencyXValues,
		YValues: mergeFrequencyYValues,
		Weeks:   weeks,
	}

	averageMergeFrequencySeriesProps := template.TimeSeriesProps{
		Series:        averageMergeFrequencySeries,
		StartEndWeeks: startEndWeek,
		InfoText:      fmt.Sprintf("AVG Merge Frequency per week: %v", util.FormatYAxisValues(averageMergeFrequencyByNWeeks)),
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

	averageReviewsSeries := template.TimeSeries{
		Title:   "Total Reviews",
		XValues: totalReviewsXValues,
		YValues: totalReviewsYValues,
		Weeks:   weeks,
	}

	averageReviewsSeriesProps := template.TimeSeriesProps{
		Series:        averageReviewsSeries,
		StartEndWeeks: startEndWeek,
		InfoText:      fmt.Sprintf("AVG Total Reviews per week: %v", util.FormatYAxisValues(averageReviewsByNWeeks)),
	}

	totalCodeChanges, averageTotalCodeChangesByNWeeks, err := store.GetTotalCodeChanges(weeks, teamMembers)

	totalCodeChangesXValues := make([]float64, len(weeks))
	totalCodeChangesYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		totalCodeChangesXValues[i] = float64(i)
		totalCodeChangesYValues[i] = float64(totalCodeChanges[week].Count)
	}
	averageCodeChangesSeries := template.TimeSeries{
		Title:   "Total Code Changes",
		XValues: totalCodeChangesXValues,
		YValues: totalCodeChangesYValues,
		Weeks:   weeks}

	averageTotalCodeChangesProps := template.TimeSeriesProps{
		Series:        averageCodeChangesSeries,
		StartEndWeeks: startEndWeek,
		InfoText:      fmt.Sprintf("AVG Total Code Changes per week: %v", util.FormatYAxisValues(averageTotalCodeChangesByNWeeks)),
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

	teamPickerProps := template.TeamPickerProps{
		Teams:        teams,
		SearchParams: url.Values{},
		SelectedTeam: team,
		BaseUrl:      "/metrics/throughput",
	}

	components := template.ThroughputMetricsPage(page, props, teamPickerProps)
	return components.Render(context.Background(), c.Response().Writer)
}

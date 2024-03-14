package handler

import (
	"fmt"

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

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	totalCommits, averageTotalCommitsByNWeeks, err := store.GetTotalCommits(weeks)

	if err != nil {
		return err
	}

	totalCommitsXValues := make([]float64, len(weeks))
	totalCommitsYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		totalCommitsXValues[i] = float64(i)
		totalCommitsYValues[i] = float64(totalCommits[week].Count)
	}

	averageTotalCommitsSeries := template.TimeSeries{
		Title:   "Total Commits",
		XValues: totalCommitsXValues,
		YValues: totalCommitsYValues,
		Weeks:   weeks,
	}

	averageTotalCommitsSeriesProps := template.TimeSeriesProps{
		Series:   averageTotalCommitsSeries,
		InfoText: fmt.Sprintf("AVG Commits per week: %v", util.FormatYAxisValues(averageTotalCommitsByNWeeks)),
	}

	totalMrsOpened, averageMrsOpenedByNWeeks, err := store.GetTotalMrsOpened(weeks)

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
		Series:   averageMrsOpenedSeries,
		InfoText: fmt.Sprintf("AVG MRs Opened per week: %v", util.FormatYAxisValues(averageMrsOpenedByNWeeks)),
	}

	mergeFrequency, averageMergeFrequencyByNWeeks, err := store.GetMergeFrequency(weeks)

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
		Series:   averageMergeFrequencySeries,
		InfoText: fmt.Sprintf("AVG Merge Frequency per week: %v", util.FormatYAxisValues(averageMergeFrequencyByNWeeks)),
	}

	totalReviews, averageReviewsByNWeeks, err := store.GetTotalReviews(weeks)

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
		Series:   averageReviewsSeries,
		InfoText: fmt.Sprintf("AVG Total Reviews per week: %v", util.FormatYAxisValues(averageReviewsByNWeeks)),
	}

	totalCodeChanges, averageTotalCodeChangesByNWeeks, err := store.GetTotalCodeChanges(weeks)

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
		Series:   averageCodeChangesSeries,
		InfoText: fmt.Sprintf("AVG Total Code Changes per week: %v", util.FormatYAxisValues(averageTotalCodeChangesByNWeeks)),
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

	components := template.ThroughputMetricsPage(page, props)
	return components.Render(context.Background(), c.Response().Writer)
}

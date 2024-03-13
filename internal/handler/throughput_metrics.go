package handler

import (
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

	totalCodeChanges, averageTotalCodeChangesByNWeeks, err := store.GetTotalCodeChanges(weeks)

	if err != nil {
		return err
	}

	totalCodeChangesXValues := make([]float64, len(weeks))
	totalCodeChangesYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		totalCodeChangesXValues[i] = float64(i)
		totalCodeChangesYValues[i] = float64(totalCodeChanges[week].Count)
	}

	props := template.ThroughputMetricsProps{
		TotalCommitsSeries:      template.TimeSeries{Title: "Total Commits", XValues: totalCommitsXValues, YValues: totalCommitsYValues, Weeks: weeks},
		AverageTotalCommits:     averageTotalCommitsByNWeeks,
		TotalMrsOpenedSeries:    template.TimeSeries{Title: "Total MRs Opened", XValues: totalMrsOpenedXValues, YValues: totalMrsOpenedYValues, Weeks: weeks},
		AverageTotalMrsOpened:   averageMrsOpenedByNWeeks,
		MergeFrequencySeries:    template.TimeSeries{Title: "Merge Frequency", XValues: mergeFrequencyXValues, YValues: mergeFrequencyYValues, Weeks: weeks},
		AverageMergeFrequency:   averageMergeFrequencyByNWeeks,
		TotalReviewsSeries:      template.TimeSeries{Title: "Total Reviews", XValues: totalReviewsXValues, YValues: totalReviewsYValues, Weeks: weeks},
		AverageTotalReviews:     averageReviewsByNWeeks,
		TotalCodeChangesSeries:  template.TimeSeries{Title: "Total Code Changes", XValues: totalCodeChangesXValues, YValues: totalCodeChangesYValues, Weeks: weeks},
		AverageTotalCodeChanges: averageTotalCodeChangesByNWeeks,
	}

	components := template.ThroughputMetricsPage(page, props)
	return components.Render(context.Background(), c.Response().Writer)
}

package handlers

import (
	"context"
	"dxta-dev/app/internal/data"
	"dxta-dev/app/internal/middlewares"
	"dxta-dev/app/internal/templates"
	"dxta-dev/app/internal/utils"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) QualityMetricsPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:     "Quality Metrics - DXTA",
		Boosted:   h.HxBoosted,
		CacheBust: a.BuildTimestamp,
		DebugMode: a.DebugMode,
	}

	tenantDatabaseUrl := r.Context().Value(middlewares.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	weeks := utils.GetLastNWeeks(time.Now(), 6*4)

	ams, err := store.GetAverageMRSize(weeks)

	if err != nil {
		return err
	}

	ard, err := store.GetAverageReviewDepth(weeks)

	if err != nil {
		return err
	}

	mmwr, err := store.GetMRsMergedWithoutReview(weeks)

	if err != nil {
		return err
	}

	amsXValues := make([]float64, len(weeks))
	amsYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		amsXValues[i] = float64(i)
		amsYValues[i] = float64(ams[week].Size)
	}

	averageMrSizeSeries := templates.TimeSeries{
		Title:   "Average Pull Request Size",
		XValues: amsXValues,
		YValues: amsYValues,
	}

	ardXValues := make([]float64, len(weeks))
	ardYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		ardXValues[i] = float64(i)
		ardYValues[i] = float64(ard[week].Depth)
	}

	averageReviewDepthSeries := templates.TimeSeries{
		Title:   "Average Review Depth",
		XValues: ardXValues,
		YValues: ardYValues,
	}

	mmwrXValues := make([]float64, len(weeks))
	mmwrYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		mmwrXValues[i] = float64(i)
		mmwrYValues[i] = float64(mmwr[week].Count)
	}

	mrsMergedWithoutReviewSeries := templates.TimeSeries{
		Title:   "Pull Requests Merged Without Review",
		XValues: mmwrXValues,
		YValues: mmwrYValues,
	}

	props := templates.QualityMetricsProps{
		AverageMrSizeSeries:          averageMrSizeSeries,
		AverageReviewDepthSeries:     averageReviewDepthSeries,
		MrsMergedWithoutReviewSeries: mrsMergedWithoutReviewSeries,
	}

	components := templates.QualityMetrics(page, props)
	return components.Render(context.Background(), c.Response().Writer)
}

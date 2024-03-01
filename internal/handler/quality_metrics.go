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

func (a *App) QualityMetricsPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &template.Page{
		Title:     "Quality Metrics - DXTA",
		Boosted:   h.HxBoosted,
		CacheBust: a.BuildTimestamp,
		DebugMode: a.DebugMode,
	}

	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	ams, amrs, err := store.GetAverageMRSize(weeks)

	if err != nil {
		return err
	}

	ard, amrrd, err := store.GetAverageReviewDepth(weeks)

	if err != nil {
		return err
	}

	mmwr, amwr, err := store.GetMRsMergedWithoutReview(weeks)

	if err != nil {
		return err
	}

	mrhm, amrh, err := store.GetAverageHandoverPerMR(weeks)

	if err != nil {
		return err
	}

	amsXValues := make([]float64, len(weeks))
	amsYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		amsXValues[i] = float64(i)
		amsYValues[i] = float64(ams[week].Size)
	}

	averageMrSizeSeries := template.TimeSeries{
		Title:   "Average Merge Request Size",
		XValues: amsXValues,
		YValues: amsYValues,
		Weeks:   weeks,
		Average: amrs,
	}

	ardXValues := make([]float64, len(weeks))
	ardYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		ardXValues[i] = float64(i)
		ardYValues[i] = float64(ard[week].Depth)
	}

	averageReviewDepthSeries := template.TimeSeries{
		Title:   "Average Review Depth",
		XValues: ardXValues,
		YValues: ardYValues,
		Weeks:   weeks,
		Average: amrrd,
	}

	ahmXValues := make([]float64, len(weeks))
	ahmYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		ahmXValues[i] = float64(i)
		ahmYValues[i] = float64(mrhm[week].Handover)
	}

	averageHandoverSeries := template.TimeSeries{
		Title:   "Average Handovers Per MR",
		XValues: ahmXValues,
		YValues: ahmYValues,
		Weeks:   weeks,
		Average: amrh,
	}

	mmwrXValues := make([]float64, len(weeks))
	mmwrYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		mmwrXValues[i] = float64(i)
		mmwrYValues[i] = float64(mmwr[week].Count)
	}

	mrsMergedWithoutReviewSeries := template.TimeSeries{
		Title:   "Pull Requests Merged Without Review",
		XValues: mmwrXValues,
		YValues: mmwrYValues,
		Weeks:   weeks,
		Average: amwr,
	}

	props := template.QualityMetricsProps{
		AverageMrSizeSeries:          averageMrSizeSeries,
		AverageReviewDepthSeries:     averageReviewDepthSeries,
		MrsMergedWithoutReviewSeries: mrsMergedWithoutReviewSeries,
		AverageHandoverTimeSeries:    averageHandoverSeries,
	}

	components := template.QualityMetricsPage(page, props)
	return components.Render(context.Background(), c.Response().Writer)
}

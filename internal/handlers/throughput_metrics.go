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

func (a *App) ThroughputMetricsPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:     "Throughput Metrics - DXTA",
		Boosted:   h.HxBoosted,
		CacheBust: a.BuildTimestamp,
		DebugMode: a.DebugMode,
	}

	tenantDatabaseUrl := r.Context().Value(middlewares.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	weeks := utils.GetLastNWeeks(time.Now(), 3*4)

	tc, err := store.GetTotalCommits(weeks)

	if err != nil {
		return err
	}

	tcXValues := make([]float64, len(weeks))
	tcYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		tcXValues[i] = float64(i)
		tcYValues[i] = float64(tc[week].Count)
	}

	tmo, err := store.GetTotalMrsOpened(weeks)

	if err != nil {
		return err
	}

	tmoXValues := make([]float64, len(weeks))
	tmoYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		tmoXValues[i] = float64(i)
		tmoYValues[i] = float64(tmo[week].Count)
	}

	mf, err := store.GetMergeFrequency(weeks)

	if err != nil {
		return err
	}

	mfXValues := make([]float64, len(weeks))
	mfYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		mfXValues[i] = float64(i)
		mfYValues[i] = float64(mf[week].Amount)
	}

	tr, err := store.GetTotalReviews(weeks)

	if err != nil {
		return err
	}

	trXValues := make([]float64, len(weeks))
	trYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		trXValues[i] = float64(i)
		trYValues[i] = float64(tr[week].Count)
	}

	tcc, err := store.GetTotalCodeChanges(weeks)

	if err != nil {
		return err
	}

	tccXValues := make([]float64, len(weeks))
	tccYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		tccXValues[i] = float64(i)
		tccYValues[i] = float64(tcc[week].Count)
	}

	props := templates.ThroughputMetricsProps{
		TotalCommitsSeries:     templates.TimeSeries{Title: "Total Commits", XValues: tcXValues, YValues: tcYValues, Weeks: weeks},
		TotalMrsOpenedSeries:   templates.TimeSeries{Title: "Total MRs Opened", XValues: tmoXValues, YValues: tmoYValues, Weeks: weeks},
		MergeFrequencySeries:   templates.TimeSeries{Title: "Merge Frequency", XValues: mfXValues, YValues: mfYValues, Weeks: weeks},
		TotalReviewsSeries:     templates.TimeSeries{Title: "Total Reviews", XValues: trXValues, YValues: trYValues, Weeks: weeks},
		TotalCodeChangesSeries: templates.TimeSeries{Title: "Total Code Changes", XValues: tccXValues, YValues: tccYValues, Weeks: weeks},
	}

	components := templates.ThroughputMetrics(page, props)
	return components.Render(context.Background(), c.Response().Writer)
}

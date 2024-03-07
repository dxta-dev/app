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

	tc, avtc, err := store.GetTotalCommits(weeks)

	if err != nil {
		return err
	}

	tcXValues := make([]float64, len(weeks))
	tcYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		tcXValues[i] = float64(i)
		tcYValues[i] = float64(tc[week].Count)
	}

	tmo, amo, err := store.GetTotalMrsOpened(weeks)

	if err != nil {
		return err
	}

	tmoXValues := make([]float64, len(weeks))
	tmoYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		tmoXValues[i] = float64(i)
		tmoYValues[i] = float64(tmo[week].Count)
	}

	mf, amf, err := store.GetMergeFrequency(weeks)

	if err != nil {
		return err
	}

	mfXValues := make([]float64, len(weeks))
	mfYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		mfXValues[i] = float64(i)
		mfYValues[i] = float64(mf[week].Amount)
	}

	tr, arx, err := store.GetTotalReviews(weeks)

	if err != nil {
		return err
	}

	trXValues := make([]float64, len(weeks))
	trYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		trXValues[i] = float64(i)
		trYValues[i] = float64(tr[week].Count)
	}

	tcc, atcc, err := store.GetTotalCodeChanges(weeks)

	if err != nil {
		return err
	}

	tccXValues := make([]float64, len(weeks))
	tccYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		tccXValues[i] = float64(i)
		tccYValues[i] = float64(tcc[week].Count)
	}

	props := template.ThroughputMetricsProps{
		TotalCommitsSeries:      template.TimeSeries{Title: "Total Commits", XValues: tcXValues, YValues: tcYValues, Weeks: weeks},
		AverageTotalCommits:     avtc,
		TotalMrsOpenedSeries:    template.TimeSeries{Title: "Total MRs Opened", XValues: tmoXValues, YValues: tmoYValues, Weeks: weeks},
		AverageTotalMrsOpened:   amo,
		MergeFrequencySeries:    template.TimeSeries{Title: "Merge Frequency", XValues: mfXValues, YValues: mfYValues, Weeks: weeks},
		AverageMergeFrequency:   amf,
		TotalReviewsSeries:      template.TimeSeries{Title: "Total Reviews", XValues: trXValues, YValues: trYValues, Weeks: weeks},
		AverageTotalReviews:     arx,
		TotalCodeChangesSeries:  template.TimeSeries{Title: "Total Code Changes", XValues: tccXValues, YValues: tccYValues, Weeks: weeks},
		AverageTotalCodeChanges: atcc,
	}

	components := template.ThroughputMetricsPage(page, props)
	return components.Render(context.Background(), c.Response().Writer)
}

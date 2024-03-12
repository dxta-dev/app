package handler

import (
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

	teams, err := store.GetTeams()

	if err != nil {
		return err
	}

	var team *data.TeamRef
	if c.Request().URL.Query().Has("team") {
		value, err := strconv.ParseInt(c.Request().URL.Query().Get("team"), 10, 64)
		if err == nil {
			team = &data.TeamRef{Id: value}
		}
	}

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	ams, amrs, err := store.GetAverageMRSize(weeks, team)

	if err != nil {
		return err
	}

	ard, amrrd, err := store.GetAverageReviewDepth(weeks, team)

	if err != nil {
		return err
	}

	mmwr, amwr, err := store.GetMRsMergedWithoutReview(weeks, team)

	if err != nil {
		return err
	}

	mrhm, amrh, err := store.GetAverageHandoverPerMR(weeks, team)

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
	}

	props := template.QualityMetricsProps{
		AverageMrSizeSeries:          averageMrSizeSeries,
		TotalAverageMrSize:           amrs,
		AverageReviewDepthSeries:     averageReviewDepthSeries,
		TotalAverageReviewDepth:      amrrd,
		MrsMergedWithoutReviewSeries: mrsMergedWithoutReviewSeries,
		TotalMrsMergedWithoutReview:  amwr,
		AverageHandoverTimeSeries:    averageHandoverSeries,
		TotalAverageHandoverTime:     amrh,
	}

	teamPickerProps := template.TeamPickerProps{
		Teams:        teams,
		SearchParams: url.Values{},
		SelectedTeam: team,
		BaseUrl:      "/metrics/quality",
	}

	components := template.QualityMetricsPage(page, props, teamPickerProps)
	return components.Render(context.Background(), c.Response().Writer)
}

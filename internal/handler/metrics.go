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

func (a *App) MetricsPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &template.Page{
		Title:     "Metrics",
		Boosted:   h.HxBoosted,
		CacheBust: a.BuildTimestamp,
		DebugMode: a.DebugMode,
	}

	var team *data.TeamRef
	if c.Request().URL.Query().Has("team") {
		value, err := strconv.ParseInt(c.Request().URL.Query().Get("team"), 10, 64)
		if err == nil {
			team = &data.TeamRef{Id: value}
		}
	}

	weeks := util.GetLastNWeeks(time.Now(), 24)

	for i, j := 0, len(weeks)-1; i < j; i, j = i+1, j-1 {
		weeks[i], weeks[j] = weeks[j], weeks[i]
	}

	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	teams, err := store.GetTeams()

	if err != nil {
		return err
	}

	averageMrSizeMap, _, err := store.GetAverageMRSize(weeks, team)

	if err != nil {
		return err
	}

	averageReviewDepthMap, _, err := store.GetAverageReviewDepth(weeks, team)

	if err != nil {
		return err
	}

	totalCommitsMap, _, err := store.GetTotalCommits(weeks, team)

	if err != nil {
		return err
	}

	totalMrsOpenedMap, _, err := store.GetTotalMrsOpened(weeks, team)

	if err != nil {
		return err
	}

	mrsMergedWithoutReviewMap, _, err := store.GetMRsMergedWithoutReview(weeks, team)

	if err != nil {
		return err
	}

	mergeFrequencyMap, _, err := store.GetMergeFrequency(weeks, team)

	if err != nil {
		return err
	}

	totalReviewsMap, _, err := store.GetTotalReviews(weeks, team)

	if err != nil {
		return err
	}

	totalCodeChanges, _, err := store.GetTotalCodeChanges(weeks, team)

	if err != nil {
		return err
	}

	metricsProps := &template.MetricsProps{
		Weeks:                 weeks,
		AverageMrSizeMap:      averageMrSizeMap,
		AverageReviewDepthMap: averageReviewDepthMap,
		MrCountMap:            mrsMergedWithoutReviewMap,
		TotalCommitsMap:       totalCommitsMap,
		TotalMrsOpenedMap:     totalMrsOpenedMap,
		TotalReviewsMap:       totalReviewsMap,
		TotalCodeChangesMap:   totalCodeChanges,
		MergeFrequencyMap:     mergeFrequencyMap,
	}

	teamPickerProps := &template.TeamPickerProps{
		Teams:        teams,
		SearchParams: url.Values{},
		SelectedTeam: team,
		BaseUrl:      "/metrics",
	}

	components := template.MetricsPage(page, *metricsProps, *teamPickerProps)
	return components.Render(context.Background(), c.Response().Writer)
}

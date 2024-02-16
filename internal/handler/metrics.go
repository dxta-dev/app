package handler

import (
	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/middleware"
	"github.com/dxta-dev/app/internal/template"
	"github.com/dxta-dev/app/internal/util"

	"time"
	"context"

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

	weeks := util.GetLastNWeeks(time.Now(), 24)

	for i, j := 0, len(weeks)-1; i < j; i, j = i+1, j-1 {
		weeks[i], weeks[j] = weeks[j], weeks[i]
	}

	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	averageMrSizeMap, err := store.GetAverageMRSize(weeks)

	if err != nil {
		return err
	}

	averageReviewDepthMap, err := store.GetAverageReviewDepth(weeks)

	if err != nil {
		return err
	}

	totalCommitsMap, err := store.GetTotalCommits(weeks)

	if err != nil {
		return err
	}

	totalMrsOpenedMap, err := store.GetTotalMrsOpened(weeks)

	if err != nil {
		return err
	}

	mrsMergedWithoutReviewMap, err := store.GetMRsMergedWithoutReview(weeks)

	if err != nil {
		return err
	}

	mergeFrequencyMap, err := store.GetMergeFrequency(weeks)

	if err != nil {
		return err
	}

	totalReviewsMap, err := store.GetTotalReviews(weeks)

	if err != nil {
		return err
	}

	totalCodeChanges, err := store.GetTotalCodeChanges(weeks)

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

	components := template.MetricsPage(page, *metricsProps)
	return components.Render(context.Background(), c.Response().Writer)
}

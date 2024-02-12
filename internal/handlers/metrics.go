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

func (a *App) Metrics(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:     "Metrics",
		Boosted:   h.HxBoosted,
		CacheBust: a.BuildTimestamp,
		DebugMode: a.DebugMode,
	}

	weeks := utils.GetLastNWeeks(time.Now(), 24)

	tenantDatabaseUrl := r.Context().Value(middlewares.TenantDatabaseURLContext).(string)

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

	metricsProps := &templates.MetricsProps{
		Weeks:                 weeks,
		AverageMrSizeMap:      averageMrSizeMap,
		AverageReviewDepthMap: averageReviewDepthMap,
		MrCountMap:            mrsMergedWithoutReviewMap,
		TotalCommitsMap:       totalCommitsMap,
		TotalMrsOpenedMap:     totalMrsOpenedMap,
	}

	components := templates.Metrics(page, *metricsProps)
	return components.Render(context.Background(), c.Response().Writer)
}

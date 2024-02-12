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

	weeks := utils.GetLastNWeeks(time.Now(), 6)

	tenantDatabaseUrl := r.Context().Value(middlewares.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	averageMrSizeData, err := store.GetAverageMRSize(weeks)

	if err != nil {
		return err
	}

	averageReviewDepthData, err := store.GetAverageReviewDepth(weeks)

	if err != nil {
		return err
	}

	mrCountData, err := store.GetMRsMergedWithoutReview(weeks)

	metricsProps := &templates.MetricsProps{
		Weeks:                  weeks,
		AverageMrSizeData:      averageMrSizeData,
		AverageReviewDepthData: averageReviewDepthData,
		MrCountData:            mrCountData,
	}

	components := templates.Metrics(page, *metricsProps)
	return components.Render(context.Background(), c.Response().Writer)
}

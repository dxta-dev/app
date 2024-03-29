package handler

import (
	"time"

	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/middleware"

	"fmt"

	"github.com/labstack/echo/v4"
)

func (a *App) GetCrawlInstancesInfo(c echo.Context) error {
	r := c.Request()
	// h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	currentTime := time.Now()

	fifteenMinutesAgo := currentTime.Add(-15 * time.Minute)
	fifteenMinutesAgoTimestamp := fifteenMinutesAgo.Unix()

	crawlInstances, err := store.GetCrawlInstances(0, fifteenMinutesAgoTimestamp)
	if err != nil {
		return err
	}

	fmt.Print("INSTANCE CRAWL-a", crawlInstances)

	return nil
}

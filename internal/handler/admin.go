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

	findGaps := data.FindGaps(time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC), fifteenMinutesAgo, crawlInstances) // This is hardcoded just for testing purposes,
	// First argument will be crawlInstances[0].Since

	if len(findGaps) > 0 {
		fmt.Println("Gaps detected:")
		for _, gap := range findGaps {
			fmt.Printf("There is a gap between %s and %s\n", gap.Since.Format("2006-01-02 15:04:05"), gap.Until.Format("2006-01-02 15:04:05"))
		}
	} else {
		fmt.Println("No gaps found.")
	}

	return nil
}

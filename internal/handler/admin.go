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

	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}

	currentTime := time.Now()

	threeMonthsAgo := currentTime.Add(-3 * time.Hour * 24 * 30)
	threeMonthsAgoTimestamp := threeMonthsAgo.Unix()

	crawlInstances, err := store.GetCrawlInstances(0, threeMonthsAgoTimestamp)
	if err != nil {
		return err
	}

	// by repositoryId

	findGaps := data.FindGaps(threeMonthsAgo, currentTime, crawlInstances)

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

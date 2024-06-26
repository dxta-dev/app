package handler

import (
	"time"

	"github.com/a-h/templ"
	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/middleware"
	admin_template "github.com/dxta-dev/app/internal/template/admin"

	"fmt"

	"github.com/labstack/echo/v4"
)

func (a *App) GetCrawlInstancesInfo(c echo.Context) error {
	r := c.Request()

	ctx := r.Context()
	store := r.Context().Value(middleware.StoreContextKey).(*data.Store)

	currentTime := time.Now()

	threeMonthsAgo := currentTime.Add(-3 * time.Hour * 24 * 30)

	crawlInstances, err := store.GetCrawlInstances(threeMonthsAgo.Unix(), currentTime.Unix())
	if err != nil {
		return err
	}

	instanceByRepo := make(map[int64]data.TimeFrameSlice)
	for _, instance := range crawlInstances {
		instanceByRepo[instance.RepositoryId] = append(instanceByRepo[instance.RepositoryId], instance.TimeFrame)
	}

	var findGaps data.TimeFrameSlice
	var components templ.Component

	for repositoryId, timeFrames := range instanceByRepo {
		findGaps = data.FindGaps(threeMonthsAgo, currentTime, timeFrames)
		if len(findGaps) > 0 {
			fmt.Printf("Gaps detected for repositoryId %d:\n", repositoryId)
			for _, gap := range findGaps {
				fmt.Printf("There is a gap between %s and %s\n", gap.Since.Format("2006-01-02 15:04:05"), gap.Until.Format("2006-01-02 15:04:05"))
			}
		} else {
			fmt.Printf("No gaps found for repositoryId %d.\n", repositoryId)
		}
		components = admin_template.DashboardAdminPage(findGaps, repositoryId, crawlInstances)
	}

	return components.Render(ctx, c.Response().Writer)
}

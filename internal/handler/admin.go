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

	// 	paramMrId := c.Param("mrid")
	// 	mrId, err := strconv.ParseInt(paramMrId, 10, 64)

	// 	if paramMrId == "" || err != nil {
	// 		return c.String(400, "")
	// 	}

	// 	parsedURL, err := url.Parse(h.HxCurrentURL)

	// 	if err != nil {
	// 		return err
	// 	}

	// 	week := parsedURL.Query().Get("week")

	// 	state := DashboardState{
	// 		week: week,
	// 		mr:   &mrId,
	// 	}

	// 	if team := parsedURL.Query().Get("team"); team != "" {
	// 		teamId, err := strconv.ParseInt(team, 10, 64)

	// 		if err != nil {
	// 			return err
	// 		}

	// 		a.State.Team = &teamId
	// 	}

	// 	nextUrl, err := getNextDashboardUrl(a, h.HxCurrentURL, state, nil, true)

	// 	if err != nil {
	// 		return err
	// 	}

	// 	c.Response().Header().Set("HX-Push-Url", nextUrl)

	// 	events, err := store.GetMergeRequestEvents(mrId)

	// 	if err != nil {
	// 		return err
	// 	}

	// 	mergeRequestInfoProps := template.MergeRequestInfoProps{
	// 		Events:         events,
	// 		DeleteEndpoint: fmt.Sprintf("/merge-request/%d", mrId),
	// 		TargetSelector: "#slide-over",
	// 	}

	// 	components := template.MergeRequestInfo(mergeRequestInfoProps)

	// 	return components.Render(context.Background(), c.Response().Writer)
	// }

	// func (a *App) RemoveMergeRequestInfo(c echo.Context) error {
	// 	r := c.Request()
	// 	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	// 	parsedURL, err := url.Parse(h.HxCurrentURL)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	week := parsedURL.Query().Get("week")

	// 	state := DashboardState{
	// 		week: week,
	// 		mr:   nil,
	// 	}

	// 	if team := parsedURL.Query().Get("team"); team != "" {
	// 		teamId, err := strconv.ParseInt(team, 10, 64)

	// 		if err != nil {
	// 			return err
	// 		}

	// 		a.State.Team = &teamId
	// 	}

	// 	nextUrl, err := getNextDashboardUrl(a, h.HxCurrentURL, state, nil, true)

	// 	if err != nil {
	// 		return err
	// 	}

	// 	c.Response().Header().Set("HX-Push-Url", nextUrl)

	// 	c.NoContent(http.StatusOK)
	return nil
}

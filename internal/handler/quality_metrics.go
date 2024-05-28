package handler

import (
	"fmt"
	"net/url"

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

	a.GenerateNonce()
	a.LoadState(r)

	currentTime := time.Now()

	threeMonthsAgo := currentTime.Add(-3 * time.Hour * 24 * 30)

	tenantDatabaseUrl := r.Context().Value(middleware.TenantDatabaseURLContext).(string)

	store := &data.Store{
		DbUrl: tenantDatabaseUrl,
	}
	crawlInstances, err := store.GetCrawlInstances(threeMonthsAgo.Unix(), currentTime.Unix())
	if err != nil {
		return err
	}

	crawlInstanceFrom := util.GetFormattedWeek(crawlInstances[0].Since)
	crawlInstanceTo := util.GetFormattedWeek(crawlInstances[len(crawlInstances)-1].Until)

	cutOffWeeks := template.CutOffWeeks{Start: crawlInstanceFrom, End: crawlInstanceTo}

	instanceByRepo := make(map[int64]data.TimeFrameSlice)
	for _, instance := range crawlInstances {
		instanceByRepo[instance.RepositoryId] = append(instanceByRepo[instance.RepositoryId], instance.TimeFrame)
	}

	var findGaps data.TimeFrameSlice

	uniqueYearWeekGaps := make(map[string]bool)

	for _, timeFrames := range instanceByRepo {
		findGaps = data.FindGaps(threeMonthsAgo, currentTime, timeFrames)
		if len(findGaps) > 0 {
			for _, gap := range findGaps {
				for d := gap.Since; d.Before(gap.Until) || d.Equal(gap.Until); d = d.AddDate(0, 0, 1) {
					year, week := d.ISOWeek()
					combination := fmt.Sprintf("%d-W%02d", year, week)

					if _, found := uniqueYearWeekGaps[combination]; !found {
						uniqueYearWeekGaps[combination] = true
					}
				}
			}
		}
	}

	teams, err := store.GetTeams()

	if err != nil {
		return err
	}

	team := a.State.Team

	teamMembers, err := store.GetTeamMembers(team)

	if err != nil {
		return err
	}

	weeks := util.GetLastNWeeks(time.Now(), 3*4)

	averageMrSize, averageMrSizeByNWeeks, err := store.GetAverageMRSize(weeks, teamMembers)

	if err != nil {
		return err
	}

	averageReviewDepth, averageReviewDepthByNWeeks, err := store.GetAverageReviewDepth(weeks, teamMembers)

	if err != nil {
		return err
	}

	mergeRequestWithoutReview, averageMrWithoutReviewByNWeeks, err := store.GetMRsMergedWithoutReview(weeks, teamMembers)

	if err != nil {
		return err
	}

	mergeRequestHandover, averageMrHandoverMetricsByNWeeks, err := store.GetAverageHandoverPerMR(weeks, teamMembers)

	if err != nil {
		return err
	}

	averageMrSizeXValues := make([]float64, len(weeks))
	averageMrSizeYValues := make([]float64, len(weeks))
	formattedAverageMrSizeYValues := make([]string, len(weeks))
	startEndWeek := make([]template.StartEndWeek, len(weeks))

	for i, week := range weeks {
		averageMrSizeXValues[i] = float64(i)
		if uniqueYearWeekGaps[week] {
			averageMrSizeYValues[i] = 0
		} else {
			averageMrSizeYValues[i] = float64(averageMrSize[week].Size)
			formattedAverageMrSizeYValues[i] = util.FormatYAxisValues(averageMrSizeYValues[i])
		}

		startWeek, endWeek, err := util.ParseYearWeek(week)
		if err != nil {
			return err
		}
		startEndWeek[i] = template.StartEndWeek{
			Start: startWeek.Format("Jan 02"),
			End:   endWeek.Format("Jan 02"),
		}
	}

	averageMrSizeSeries := template.TimeSeries{
		Title:   "Average Merge Request Size",
		XValues: averageMrSizeXValues,
		YValues: averageMrSizeYValues,
		Weeks:   weeks,
	}

	averageMrSizeSeriesProps := template.TimeSeriesProps{
		Series:           averageMrSizeSeries,
		StartEndWeeks:    startEndWeek,
		CutOffWeeks:      cutOffWeeks,
		FormattedYValues: formattedAverageMrSizeYValues,
		InfoText:         fmt.Sprintf("AVG Size per week: %v", util.FormatYAxisValues(averageMrSizeByNWeeks)),
	}

	averageReviewDepthXValues := make([]float64, len(weeks))
	averageReviewDepthYValues := make([]float64, len(weeks))
	formattedAverageReviewDepthYValues := make([]string, len(weeks))

	for i, week := range weeks {
		averageReviewDepthXValues[i] = float64(i)
		if uniqueYearWeekGaps[week] {
			averageReviewDepthYValues[i] = 0
		} else {
			averageReviewDepthYValues[i] = float64(averageReviewDepth[week].Depth)
			formattedAverageReviewDepthYValues[i] = util.FormatYAxisValues(averageReviewDepthYValues[i])
		}
	}

	averageReviewDepthSeries := template.TimeSeries{
		Title:   "Average Review Depth",
		XValues: averageReviewDepthXValues,
		YValues: averageReviewDepthYValues,
		Weeks:   weeks,
	}

	averageReviewDepthSeriesProps := template.TimeSeriesProps{
		Series:           averageReviewDepthSeries,
		StartEndWeeks:    startEndWeek,
		CutOffWeeks:      cutOffWeeks,
		FormattedYValues: formattedAverageReviewDepthYValues,
		InfoText:         fmt.Sprintf("AVG Depth per week: %v", util.FormatYAxisValues(averageReviewDepthByNWeeks)),
	}

	averageMrHandoverMetricsByNWeeksXValues := make([]float64, len(weeks))
	averageMrHandoverMetricsByNWeeksYValues := make([]float64, len(weeks))
	formattedAverageMrHandoverMetricsByNWeeksYValues := make([]string, len(weeks))

	for i, week := range weeks {
		averageMrHandoverMetricsByNWeeksXValues[i] = float64(i)
		if uniqueYearWeekGaps[week] {
			averageMrHandoverMetricsByNWeeksYValues[i] = 0
		} else {
			averageMrHandoverMetricsByNWeeksYValues[i] = float64(mergeRequestHandover[week].Handover)
			formattedAverageMrHandoverMetricsByNWeeksYValues[i] = util.FormatYAxisValues(averageMrHandoverMetricsByNWeeksYValues[i])
		}
	}

	averageHandoverSeries := template.TimeSeries{
		Title:   "Average Handovers Per MR",
		XValues: averageMrHandoverMetricsByNWeeksXValues,
		YValues: averageMrHandoverMetricsByNWeeksYValues,
		Weeks:   weeks,
	}

	averageHandoverSeriesProps := template.TimeSeriesProps{
		Series:           averageHandoverSeries,
		StartEndWeeks:    startEndWeek,
		CutOffWeeks:      cutOffWeeks,
		FormattedYValues: formattedAverageMrHandoverMetricsByNWeeksYValues,
		InfoText:         fmt.Sprintf("AVG Handovers per week: %v", util.FormatYAxisValues(averageMrHandoverMetricsByNWeeks)),
	}

	mergeRequestWithoutReviewXValues := make([]float64, len(weeks))
	mergeRequestWithoutReviewYValues := make([]float64, len(weeks))
	formattedMergeRequestWithoutReviewYValues := make([]string, len(weeks))

	for i, week := range weeks {
		mergeRequestWithoutReviewXValues[i] = float64(i)
		if uniqueYearWeekGaps[week] {
			mergeRequestWithoutReviewYValues[i] = 0
		} else {
			mergeRequestWithoutReviewYValues[i] = float64(mergeRequestWithoutReview[week].Count)
			formattedMergeRequestWithoutReviewYValues[i] = util.FormatYAxisValues(mergeRequestWithoutReviewYValues[i])
		}
	}

	mrsMergedWithoutReviewSeries := template.TimeSeries{
		Title:   "Pull Requests Merged Without Review",
		XValues: mergeRequestWithoutReviewXValues,
		YValues: mergeRequestWithoutReviewYValues,
		Weeks:   weeks,
	}

	mrsMergedWithoutReviewSeriesProps := template.TimeSeriesProps{
		Series:           mrsMergedWithoutReviewSeries,
		StartEndWeeks:    startEndWeek,
		CutOffWeeks:      cutOffWeeks,
		FormattedYValues: formattedMergeRequestWithoutReviewYValues,
		InfoText:         fmt.Sprintf("Total Merged without Review: %v", util.FormatYAxisValues(averageMrWithoutReviewByNWeeks)),
	}

	props := template.QualityMetricsProps{
		AverageMrSizeSeriesProps:          averageMrSizeSeriesProps,
		AverageReviewDepthSeriesProps:     averageReviewDepthSeriesProps,
		MrsMergedWithoutReviewSeriesProps: mrsMergedWithoutReviewSeriesProps,
		AverageHandoverTimeSeriesProps:    averageHandoverSeriesProps,
	}

	var templTeams []template.Team

	for _, team := range teams {
		params := url.Values{}
		params.Set("team", fmt.Sprint(team.Id))
		teamUrl, err := a.GetUrlAppState(r.URL.Path, params)
		if err != nil {
			return err
		}
		templTeams = append(templTeams, template.Team{
			Id:   team.Id,
			Name: team.Name,
			Url:  teamUrl,
		})
	}

	teamPickerProps := template.TeamPickerProps{
		Teams:        templTeams,
		SelectedTeam: team,
		NoTeamUrl:    r.URL.Path,
	}

	navState, err := a.GetNavState()

	if err != nil {
		return err
	}

	page := &template.Page{
		Title:     "Quality Metrics - DXTA",
		Boosted:   h.HxBoosted,
		CacheBust: a.BuildTimestamp,
		DebugMode: a.DebugMode,
		NavState:  navState,
		Nonce:     a.Nonce,
	}

	components := template.QualityMetricsPage(page, props, teamPickerProps)
	return components.Render(context.Background(), c.Response().Writer)
}

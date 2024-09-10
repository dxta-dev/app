package handler

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/dxta-dev/app/internal/data"
	"github.com/dxta-dev/app/internal/middleware"
	"github.com/dxta-dev/app/internal/template"
	"github.com/dxta-dev/app/internal/util"

	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) QualityMetricsPage(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	store := r.Context().Value(middleware.StoreContextKey).(*data.Store)

	a.GenerateNonce()
	a.LoadState(r)

	ctx := r.Context()

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

	var wg sync.WaitGroup

	errCh := make(chan error, 4)

	var (
		averageMrSize                    map[string]data.AverageMRSizeByWeek
		averageMrSizeByNWeeks            float64
		averageReviewDepth               map[string]data.AverageMrReviewDepthByWeek
		averageReviewDepthByNWeeks       float64
		mergeRequestWithoutReview        map[string]data.MrCountByWeek
		averageMrWithoutReviewByNWeeks   float64
		mergeRequestHandover             map[string]data.AverageHandoverPerMR
		averageMrHandoverMetricsByNWeeks float64
	)

	wg.Add(4)

	go func() {
		defer wg.Done()
		averageMrSize, averageMrSizeByNWeeks, err = store.GetAverageMRSize(weeks, teamMembers)
		if err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		averageReviewDepth, averageReviewDepthByNWeeks, err = store.GetAverageReviewDepth(weeks, teamMembers)
		if err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		mergeRequestWithoutReview, averageMrWithoutReviewByNWeeks, err = store.GetMRsMergedWithoutReview(weeks, teamMembers)
		if err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		mergeRequestHandover, averageMrHandoverMetricsByNWeeks, err = store.GetAverageHandoverPerMR(weeks, teamMembers)
		if err != nil {
			errCh <- err
			return
		}
	}()

	wg.Wait()

	close(errCh)
	for e := range errCh {
		if e != nil {
			err = e
			break
		}
	}

	if err != nil {
		return err
	}

	averageLifecycle, averageLifecycleByNWeeks, err := store.GetAverageLifecycleDuration(weeks, teamMembers)

	if err != nil {
		return err
	}

	averageMrSizeXValues := make([]float64, len(weeks))
	averageMrSizeYValues := make([]float64, len(weeks))
	startEndWeek := make([]template.StartEndWeek, len(weeks))

	for i, week := range weeks {
		averageMrSizeXValues[i] = float64(i)
		averageMrSizeYValues[i] = float64(averageMrSize[week].Size)
		startWeek, endWeek, err := util.ParseYearWeek(week)
		if err != nil {
			return err
		}
		startEndWeek[i] = template.StartEndWeek{
			Start: startWeek.Format("Jan 02"),
			End:   endWeek.Format("Jan 02"),
		}
	}

	formattedAverageMrSizeYValues := make([]string, len(averageMrSizeYValues))

	for i, value := range averageMrSizeYValues {
		formattedAverageMrSizeYValues[i] = util.FormatYAxisValues(value)
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
		FormattedYValues: formattedAverageMrSizeYValues,
		InfoText:         fmt.Sprintf("AVG Size per week: %v", util.FormatYAxisValues(averageMrSizeByNWeeks)),
	}

	averageReviewDepthXValues := make([]float64, len(weeks))
	averageReviewDepthYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		averageReviewDepthXValues[i] = float64(i)
		averageReviewDepthYValues[i] = float64(averageReviewDepth[week].Depth)
	}

	formattedAverageReviewDepthYValues := make([]string, len(averageReviewDepthYValues))

	for i, value := range averageReviewDepthYValues {
		formattedAverageReviewDepthYValues[i] = util.FormatYAxisValues(value)
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
		FormattedYValues: formattedAverageReviewDepthYValues,
		InfoText:         fmt.Sprintf("AVG Depth per week: %v", util.FormatYAxisValues(averageReviewDepthByNWeeks)),
	}

	averageMrHandoverMetricsByNWeeksXValues := make([]float64, len(weeks))
	averageMrHandoverMetricsByNWeeksYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		averageMrHandoverMetricsByNWeeksXValues[i] = float64(i)
		averageMrHandoverMetricsByNWeeksYValues[i] = float64(mergeRequestHandover[week].Handover)
	}

	formattedAverageMrHandoverMetricsByNWeeksYValues := make([]string, len(averageMrHandoverMetricsByNWeeksYValues))

	for i, value := range averageMrHandoverMetricsByNWeeksYValues {
		formattedAverageMrHandoverMetricsByNWeeksYValues[i] = util.FormatYAxisValues(value)
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
		FormattedYValues: formattedAverageMrHandoverMetricsByNWeeksYValues,
		InfoText:         fmt.Sprintf("AVG Handovers per week: %v", util.FormatYAxisValues(averageMrHandoverMetricsByNWeeks)),
	}

	mergeRequestWithoutReviewXValues := make([]float64, len(weeks))
	mergeRequestWithoutReviewYValues := make([]float64, len(weeks))

	for i, week := range weeks {
		mergeRequestWithoutReviewXValues[i] = float64(i)
		mergeRequestWithoutReviewYValues[i] = float64(mergeRequestWithoutReview[week].Count)
	}

	formattedMergeRequestWithoutReviewYValues := make([]string, len(mergeRequestWithoutReviewYValues))

	for i, value := range mergeRequestWithoutReviewYValues {
		formattedMergeRequestWithoutReviewYValues[i] = util.FormatYAxisValues(value)
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
		FormattedYValues: formattedMergeRequestWithoutReviewYValues,
		InfoText:         fmt.Sprintf("Total Merged without Review: %v", util.FormatYAxisValues(averageMrWithoutReviewByNWeeks)),
	}

	averageLifecycleXValues := make([]float64, len(weeks))
	averageLifecycleYValues := make([]float64, len(weeks))
	normalizedAverageLifecycleYValues := make([]float64, len(weeks))
	maxAverageLifecycleYValue := 0.0

	for i, week := range weeks {
		averageLifecycleXValues[i] = float64(i)
		averageLifecycleYValues[i] = float64(averageLifecycle[week].Lifecycle)
		if maxAverageLifecycleYValue < averageLifecycle[week].Lifecycle {
			maxAverageLifecycleYValue = averageLifecycle[week].Lifecycle
		}
	}

	averageLifecycleYValueNormalFactor := util.GetTimeAxisNormalizationFactor(maxAverageLifecycleYValue)
	for i, _ := range averageLifecycleYValues {
		normalizedAverageLifecycleYValues[i] = averageLifecycleYValues[i] * averageLifecycleYValueNormalFactor
	}

	formattedAverageLifecycleYValues := make([]string, len(averageLifecycleYValues))

	for i, value := range averageLifecycleYValues {
		formattedAverageLifecycleYValues[i] = util.FormatTimeAxisValue(value)
	}

	averageLifecycleSeries := template.TimeSeries{
		Title:   "Average Lifecyle Duration",
		XValues: averageLifecycleXValues,
		YValues: normalizedAverageLifecycleYValues,
		Weeks:   weeks,
	}

	averageLifecycleSeriesProps := template.TimeSeriesProps{
		Series:           averageLifecycleSeries,
		StartEndWeeks:    startEndWeek,
		FormattedYValues: formattedAverageLifecycleYValues,
		InfoText:         fmt.Sprintf("AVG Lifecycle per week: %v", util.FormatTimeAxisValue(averageLifecycleByNWeeks)),
	}

	props := template.QualityMetricsProps{
		AverageMrSizeSeriesProps:          averageMrSizeSeriesProps,
		AverageReviewDepthSeriesProps:     averageReviewDepthSeriesProps,
		MrsMergedWithoutReviewSeriesProps: mrsMergedWithoutReviewSeriesProps,
		AverageHandoverTimeSeriesProps:    averageHandoverSeriesProps,
		AverageLifecycleSeriesProps:       averageLifecycleSeriesProps,
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
		RouteId:   "/metrics/quality",
		Title:     "Quality Metrics - DXTA",
		Boosted:   h.HxBoosted,
		CacheBust: a.BuildTimestamp,
		DebugMode: a.DebugMode,
		NavState:  navState,
		Nonce:     a.Nonce,
	}

	components := template.QualityMetricsPage(page, props, teamPickerProps)
	return components.Render(ctx, c.Response().Writer)
}

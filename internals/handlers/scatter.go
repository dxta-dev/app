package handlers

import (
	"context"
	"dxta-dev/app/internals/templates"
	"strconv"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func parseInt(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}

func parseFloat64(str string) float64 {
	v, _ := strconv.ParseFloat(str, 64)
	return v
}

func readData() ([]float64, []float64) {
	var xvalues []float64
	var yvalues []float64

	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1).Truncate(24 * time.Hour)

	times := []time.Time{
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 16, 51, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 18, 5, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 18, 10, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 18, 15, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 18, 20, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 18, 25, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 18, 30, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 18, 35, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 18, 40, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 18, 45, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day()+1, 19, 0, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day()+2, 20, 0, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day()+3, 23, 0, 0, 0, time.UTC),
	}

	for _, time := range times {
		xSecondsValue := float64(time.Unix() - startOfWeek.Unix())
		xvalues = append(xvalues, xSecondsValue)
		yvalues = append(yvalues, 60*60*24)
	}

	return xvalues, yvalues
}

func (a *App) Scatter(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	page := &templates.Page{
		Title:   "Charts",
		Boosted: h.HxBoosted,
	}

	var chartData []templates.ScatterSeries
	xValues, yValues := readData()

	chartData = append(chartData, templates.ScatterSeries{
		Title:   "Random Series",
		XValues: xValues,
		YValues: yValues,
	})


	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1).Truncate(24 * time.Hour)

	components := templates.Scatter(page, chartData, startOfWeek)
	return components.Render(context.Background(), c.Response().Writer)
}

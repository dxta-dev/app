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

	// Hardcoded dates
	times := []time.Time{
		time.Date(2023, 12, 19, 16, 51, 0, 0, time.UTC),
		time.Date(2023, 12, 19, 18, 1, 0, 0, time.UTC),
		time.Date(2023, 12, 22, 23, 0, 0, 0, time.UTC),
	}

	for _, time := range times {
		xSecondsValue := float64(time.Unix() - startOfWeek.Unix())
		xvalues = append(xvalues, xSecondsValue)
		yvalues = append(yvalues, 50)
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

	components := templates.Scatter(page, chartData)
	return components.Render(context.Background(), c.Response().Writer)
}

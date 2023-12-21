package handlers

import (
	"context"
	"dxta-dev/app/internals/templates"
	"strconv"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	"github.com/wcharczuk/go-chart/v2"
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

	// Hardcoded dates
	dates := []string{
		"2023,12,19,16,51.1948264984227130",
		"2023,12,19,18,1.7940833333333333",
		"2023,12,22,18,10.0383889931207000",
	}

	for _, dateStr := range dates {
		parts := chart.SplitCSV(dateStr)
		year := parseInt(parts[0])
		month := parseInt(parts[1])
		day := parseInt(parts[2])
		hour := parseInt(parts[3])
		// elapsedMillis := parseFloat64(parts[4]) -> we will later implement Yvalue so this will stay commented out
		timeValue := time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC)
		xmilisecValue := float64(timeValue.UnixNano()) / 1e6
		xvalues = append(xvalues, xmilisecValue)
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

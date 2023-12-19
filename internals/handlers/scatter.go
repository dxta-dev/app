package handlers

import (
	"context"
	"dxta-dev/app/internals/templates"
	"fmt"
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
	var xmilisecValue float64
	err := chart.ReadLines("internals/handlers/request.csv", func(line string) error {
		parts := chart.SplitCSV(line)
		year := parseInt(parts[0])
		month := parseInt(parts[1])
		day := parseInt(parts[2])
		hour := parseInt(parts[3])
		// elapsedMillis := parseFloat64(parts[4]) -> we will later implement Yvalue so this will stay commented out
		timeValue := time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC)
		xmilisecValue = float64(timeValue.UnixNano()) / 1e6
		xvalues = append(xvalues, xmilisecValue)
		yvalues = append(yvalues, 50)
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
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

package handlers

import (
	"context"
	"dxta-dev/app/internals/templates"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	"github.com/wcharczuk/go-chart/v2"
)

func (a *App) Scatter(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	page := &templates.Page{
		Title:   "Charts",
		Boosted: h.HxBoosted,
	}

	var chartData []templates.ScatterSeries

	chartData = append(chartData, templates.ScatterSeries{
		Title:   "Random Series",
		XValues: chart.Seq{Sequence: chart.NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values(),
		YValues: chart.Seq{Sequence: chart.NewRandomSequence().WithLen(100).WithMin(0).WithMax(100)}.Values(),
	})

	chartData = append(chartData, templates.ScatterSeries{
		Title:   "Random Series 2",
		XValues: chart.Seq{Sequence: chart.NewLinearSequence().WithStart(50).WithEnd(55)}.Values(),
		YValues: chart.Seq{Sequence: chart.NewRandomSequence().WithLen(10).WithMin(0).WithMax(10)}.Values(),
	})

	chartData = append(chartData, templates.ScatterSeries{
		Title:   "Series",
		XValues: []float64{1, 2, 3, 4, 5},
		YValues: []float64{1, 2, 3, 4, 5},
	})

	components := templates.Scatter(page, chartData)
	return components.Render(context.Background(), c.Response().Writer)
}

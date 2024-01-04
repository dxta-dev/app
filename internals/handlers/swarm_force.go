package handlers

import (
	"context"
	"dxta-dev/app/internals/templates"
	"dxta-dev/app/internals/data"
	"dxta-dev/app/internals/graphs"
	"time"
	"sort"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)


func series1() templates.ScatterSeries {
	var xvalues []float64
	var yvalues []float64

	startOfWeek := time.Unix(1696204800, 0)

	var times []time.Time

	sort.Sort(data.DataList)

	for _, d := range data.DataList {
		t := time.Unix(d.Timestamp/1000, 0)
		times = append(times, t)
	}

	for _, time := range times {
		xSecondsValue := float64(time.Unix() - startOfWeek.Unix())
		xvalues = append(xvalues, xSecondsValue)
		yvalues = append(yvalues, 60*60*12)
	}

	xvalues, yvalues = graphs.Beehive(xvalues, yvalues)





	return templates.ScatterSeries{
		Title:   "series 1",
		XValues: xvalues,
		YValues: yvalues,
	}
}

func (a *App) ScatterForce(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	page := &templates.Page{
		Title:   "Charts",
		Boosted: h.HxBoosted,
	}

	var chartData []templates.ScatterSeries

	chartData = append(chartData, series1())


	startOfWeek := time.Unix(1696204800, 0)

	components := templates.Scatter(page, chartData, startOfWeek)
	return components.Render(context.Background(), c.Response().Writer)
}

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
	"github.com/wcharczuk/go-chart/v2/drawing"
)


func series1() templates.SwarmSeries {
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

	xvalues, yvalues = graphs.Beehive(xvalues, yvalues, 1400, 200, 5)

	colors := []drawing.Color{}

	for i := 0; i < len(xvalues); i++ {
		switch data.DataList[i].EventType {
		case data.Commited:
			colors = append(colors, drawing.ColorBlue)
		case data.Reviewed:
			colors = append(colors, drawing.ColorGreen)
		case data.Merged:
			colors = append(colors, drawing.ColorRed)
		default:
			colors = append(colors, drawing.ColorGreen)
		}
	}

	return templates.SwarmSeries{
		Title:   "series 1",
		DotColors: colors,
		XValues: xvalues,
		YValues: yvalues,
	}
}

func (a *App) Swarm(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	page := &templates.Page{
		Title:   "Charts",
		Boosted: h.HxBoosted,
	}

	var chartData []templates.SwarmSeries

	chartData = append(chartData, series1())


	startOfWeek := time.Unix(1696204800, 0)

	components := templates.Swarm(page, chartData, startOfWeek)
	return components.Render(context.Background(), c.Response().Writer)
}

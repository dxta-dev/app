package handlers

import (
	"context"
	"dxta-dev/app/internals/templates"
	"dxta-dev/app/internals/graphs"
	"dxta-dev/app/internals/data"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)


func series1() templates.ScatterSeries {
	var xvalues []float64
	var yvalues []float64

	startOfWeek := time.Unix(1696204800, 0)

	var times []time.Time

	for _, d := range data.DataList {
		t := time.Unix(d.Timestamp/1000, 0)
		times = append(times, t)
	}


	for _, time := range times {
		xSecondsValue := float64(time.Unix() - startOfWeek.Unix())
		xvalues = append(xvalues, xSecondsValue)
		yvalues = append(yvalues, 60*60*12)
	}

	graph := graphs.NewGraph()

	for i := 0; i < len(times); i++ {
		graph.AddNode(graphs.NewNode(xvalues[i], yvalues[i]))
	}

	graphs.ForceDirectedGraphLayout(graph, 1000)

	for i := 0; i < len(graph.Nodes); i++ {
		xvalues[i] = graph.Nodes[i].Position.X
		yvalues[i] = graph.Nodes[i].Position.Y
	}

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

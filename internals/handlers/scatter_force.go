package handlers

import (
	"context"
	"dxta-dev/app/internals/templates"
	"dxta-dev/app/internals/graphs"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func series1() templates.ScatterSeries {
	var xvalues []float64
	var yvalues []float64

	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1).Truncate(24 * time.Hour)

	times := []time.Time{
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 18, 5, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 18, 10, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 18, 15, 0, 0, time.UTC),
	}

	for _, time := range times {
		xSecondsValue := float64(time.Unix() - startOfWeek.Unix())
		xvalues = append(xvalues, xSecondsValue)
		yvalues = append(yvalues, 60*60*24)
	}

	graph := graphs.NewGraph()

	graph.AddNode(graphs.NewNode(xvalues[0], yvalues[0]))
	graph.AddNode(graphs.NewNode(xvalues[1], yvalues[1]))
	graph.AddNode(graphs.NewNode(xvalues[2], yvalues[2]))
	graph.AddEdge(0, 1)
	graph.AddEdge(0, 2)
	graph.AddEdge(1, 2)

	graphs.ForceDirectedGraphLayout(graph, 1)

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

	components := templates.Scatter(page, chartData)
	return components.Render(context.Background(), c.Response().Writer)
}

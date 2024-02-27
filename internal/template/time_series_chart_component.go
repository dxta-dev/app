package template

import (
	"bytes"
	"context"
	"io"
	"log"
	"strings"

	"github.com/a-h/templ"
	"github.com/wcharczuk/go-chart/v2"
)

type TimeSeries struct {
	Title   string
	XValues []float64
	YValues []float64
	Weeks   []string
}

func TimeSeriesChart(series TimeSeries) templ.Component {
	mainSeries := chart.ContinuousSeries{
		Style: chart.Style{
			StrokeWidth: chart.Disabled,
			DotWidth:    5,
			DotColor:    chart.ColorBlue,
		},
		Name:    series.Title,
		XValues: series.XValues,
		YValues: series.YValues,
	}

	lineSeries := chart.ContinuousSeries{
		Name:    series.Title,
		XValues: series.XValues,
		YValues: series.YValues,
	}

	graph := chart.Chart{
		Series: []chart.Series{
			lineSeries,
			mainSeries,
		},
		Height: 300,
		Width:  650,
	}

	for i, week := range series.Weeks {
		week_part := strings.Split(week, "-")
		graph.XAxis.Ticks = append(graph.XAxis.Ticks, chart.Tick{
			Value: float64(i),
			Label: week_part[1],
		})
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.SVG, buffer)

	if err != nil {
		log.Fatal(err)
	}

	html := buffer.String()

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, html)
		return err
	})
}

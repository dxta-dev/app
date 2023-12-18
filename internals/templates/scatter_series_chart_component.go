package templates

import (
	"bytes"
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/wcharczuk/go-chart/v2"
)

type DotType int

const (
	DotTypeCircle DotType = iota
	DotTypeSquare
	DotTypeDiamond
)

type ScatterSeries struct {
	Title    string
	DotTypes []DotType
	XValues  []float64
	YValues  []float64
}

func ScatterSeriesChart(series ScatterSeries) templ.Component {
	mainSeries := chart.ContinuousSeries{
		Style: chart.Style{
			StrokeWidth: chart.Disabled,
			DotWidth:    5,
		},
		Name:    series.Title,
		XValues: series.XValues,
		YValues: series.YValues,
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			GridMajorStyle: chart.Style{
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 1.0,
			},
			Ticks: []chart.Tick{
				{Value: 1, Label: "Monday"},
				{Value: 2, Label: "Tuesday"},
				{Value: 3, Label: "Wednesday"},
				{Value: 4, Label: "Thursday"},
				{Value: 5, Label: "Friday"},
				{Value: 6, Label: "Saturday"},
				{Value: 7, Label: "Sunday"},
			},
		},
		Series: []chart.Series{
			mainSeries,
		},
	}

	for _, tick := range graph.XAxis.Ticks {
		gridLine := chart.ContinuousSeries{
			XValues: []float64{tick.Value, tick.Value},
			YValues: []float64{0, 100},
			Style: chart.Style{
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 1.0,
			},
		}
		graph.Series = append(graph.Series, gridLine)
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.SVG, buffer)

	if err != nil {
		panic(err)
	}

	html := buffer.String()

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, html)
		return err
	})
}

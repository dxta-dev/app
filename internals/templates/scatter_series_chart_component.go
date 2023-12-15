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
		Series: []chart.Series{
			mainSeries,
		},
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

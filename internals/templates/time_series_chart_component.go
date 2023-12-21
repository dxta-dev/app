package templates

import (
	"bytes"
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/wcharczuk/go-chart/v2"
)

type TimeSeries struct {
	Title   string
	XValues []float64
	YValues []float64
}

func TimeSeriesChart(series TimeSeries) templ.Component {
	mainSeries := chart.ContinuousSeries{
		Name:    series.Title,
		XValues: series.XValues,
		YValues: series.YValues,
	}

	smaSeries := &chart.SMASeries{
		InnerSeries: mainSeries,
	}

	graph := chart.Chart{
		Series: []chart.Series{
			mainSeries,
			smaSeries,
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

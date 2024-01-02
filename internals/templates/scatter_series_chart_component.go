package templates

import (
	"bytes"
	"context"
	"io"
	"time"

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

func ScatterSeriesChart(series ScatterSeries, startOfWeek time.Time) templ.Component {
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
		},
		YAxis: chart.YAxis{
			Style: chart.Hidden(),
		},
		Series: []chart.Series{
			mainSeries,
		},
		Width: 1425,
		Height: 227,
	}

	startOfWeekSeconds := startOfWeek.Unix()
	for i := 0; i < 7; i++ {
		secondsFromStartOfWeek := startOfWeekSeconds + int64(i*24*60*60)
		secondsForEachDay := int64(i * 24 * 60 * 60)
		dateLabel := time.Unix(secondsFromStartOfWeek, 0).Format("Mon 02")
		graph.XAxis.Ticks = append(graph.XAxis.Ticks, chart.Tick{
			Value: float64(secondsForEachDay),
			Label: dateLabel,
		})
	}

	for _, tick := range graph.XAxis.Ticks {
		gridLine := chart.ContinuousSeries{
			XValues: []float64{tick.Value, tick.Value},
			YValues: []float64{0, 60*60*24},
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

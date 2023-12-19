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

func ScatterSeriesChart(series ScatterSeries) templ.Component {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

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
		Series: []chart.Series{
			mainSeries,
		},
	}

	for i := startOfWeek.Day(); i <= endOfWeek.Day(); i++ {
		day := time.Date(now.Year(), now.Month(), i, 0, 0, 0, 0, time.UTC)
		dateLabel := day.Format("Mon 02")
		graph.XAxis.Ticks = append(graph.XAxis.Ticks, chart.Tick{Value: float64(i), Label: dateLabel})
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

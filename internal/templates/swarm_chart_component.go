package templates

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/a-h/templ"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type SwarmSeries struct {
	Title     string
	DotColors []drawing.Color
	XValues   []float64
	YValues   []float64
}

func SwarmChart(series SwarmSeries, startOfWeek time.Time) templ.Component {

	colorProvider := func(xr, yr chart.Range, index int, x, y float64) drawing.Color {
		if len(series.DotColors) > index {
			return series.DotColors[index]
		}
		return chart.ColorOrange
	}

	mainSeries := chart.ContinuousSeries{
		Style: chart.Style{
			StrokeWidth:      chart.Disabled,
			DotWidth:         5,
			DotColorProvider: colorProvider,
		},
		Name:    series.Title,
		XValues: series.XValues,
		YValues: series.YValues,
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
			Style: chart.Style{
				StrokeColor: chart.ColorBlack,
			},
			GridMajorStyle: chart.Hidden(),
			GridMinorStyle: chart.Hidden(),
		},
		YAxis: chart.YAxis{
			Style: chart.Hidden(),
			GridMajorStyle: chart.Hidden(),
			GridMinorStyle: chart.Hidden(),
		},
		YAxisSecondary: chart.YAxis{
			Style: chart.Hidden(),
			GridMajorStyle: chart.Hidden(),
			GridMinorStyle: chart.Hidden(),
		},
		Series: []chart.Series{
			mainSeries,
		},
		Width:  1405,
		Height: 227,
	}

	startOfWeekSeconds := startOfWeek.Unix()
	for i := 0; i < 8; i++ {
		secondsFromStartOfWeek := startOfWeekSeconds + int64(i*24*60*60)
		secondsForEachDay := int64(i * 24 * 60 * 60)
		dateLabel := time.Unix(secondsFromStartOfWeek-24*60*60, 0).Format("Mon 02")
		graph.XAxis.Ticks = append(graph.XAxis.Ticks, chart.Tick{
			Value: float64(secondsForEachDay),
			Label: dateLabel,
		})
	}

	for _, tick := range graph.XAxis.Ticks {
		gridLine := chart.ContinuousSeries{
			XValues: []float64{tick.Value, tick.Value},
			YValues: []float64{0, 24*60*60},
			Style: chart.Style{
				StrokeWidth: 1.0,
				StrokeColor: chart.ColorBlack,
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

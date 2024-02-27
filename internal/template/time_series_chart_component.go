package template

import (
	"bytes"
	"context"
	"io"
	"log"
	"math"
	"strconv"
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

func getYAxisValues(yValues []float64) []float64 {
	if len(yValues) == 0 {
		return []float64{0, 1, 2, 3, 4}
	}

	highest := yValues[0]

	for _, num := range yValues {
		if num > highest {
			highest = num
		}
	}

	if highest < 1 {
		return []float64{0, 0.25, 0.5, 0.75, 1}
	}

	if highest < 1.6 {
		return []float64{0, 0.4, 0.8, 1.2, 1.6}
	}

	if highest < 2 {
		return []float64{0, 0.5, 1, 1.5, 2}
	}

	if highest < 3.2 {
		return []float64{0, 0.8, 1.6, 2.4, 3.2}
	}

	if highest <  4 {
		return []float64{0, 1, 2, 3, 4}
	}

	if highest < 4.8 {
		return []float64{0, 1.2, 2.4, 3.6, 4.8}
	}

	if highest < 6 {
		return []float64{0, 1.5, 3, 4.5, 6}
	}

	if highest < 7.2 {
		return []float64{0, 1.8, 3.6, 5.4, 7.2}
	}

	if highest < 8 {
		return []float64{0, 2, 4, 6, 8}
	}

	if highest < 12 {
		return []float64{0, 3, 6, 9, 12}
	}

	if highest < 16 {
		return []float64{0, 4, 8, 12, 16}
	}

	if highest < 20 {
		return []float64{0, 5, 10, 15, 20}
	}

	if highest < 24 {
		return []float64{0, 6, 12, 18, 24}
	}

	if highest < 28 {
		return []float64{0, 7, 14, 21, 28}
	}

	if highest < 32 {
		return []float64{0, 8, 16, 24, 32}
	}

	if highest < 40 {
		return []float64{0, 10, 20, 30, 40}
	}

	if highest < 64 {
		return []float64{0, 16, 32, 48, 64}
	}

	if highest < 80 {
		return []float64{0, 20, 40, 60, 80}
	}

	if highest < 96 {
		return []float64{0, 24, 48, 72, 96}
	}

	highest = float64(int(math.Ceil(highest/20) * 20))

	lowest := 0.0
	percent25 := highest * 0.25
	percent50 := highest * 0.5
	percent75 := highest * 0.75

	return []float64{lowest, percent25, percent50, percent75, highest}
}

func TimeSeriesChart(series TimeSeries) templ.Component {
	YAxisValues := getYAxisValues(series.YValues)

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
		YAxis: chart.YAxis{
			Ticks: []chart.Tick{
				{Value: YAxisValues[0], Label: strconv.FormatFloat(YAxisValues[0], 'f', -1, 64)},
				{Value: YAxisValues[1], Label: strconv.FormatFloat(YAxisValues[1], 'f', -1, 64)},
				{Value: YAxisValues[2], Label: strconv.FormatFloat(YAxisValues[2], 'f', -1, 64)},
				{Value: YAxisValues[3], Label: strconv.FormatFloat(YAxisValues[3], 'f', -1, 64)},
				{Value: YAxisValues[4], Label: strconv.FormatFloat(YAxisValues[4], 'f', -1, 64)},
			},
		},
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

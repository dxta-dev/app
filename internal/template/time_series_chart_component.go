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

	if highest < 5 {
		return []float64{0, 1, 2, 3, 4}
	}

	highest = float64(int(math.Round(highest/10) * 10))

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

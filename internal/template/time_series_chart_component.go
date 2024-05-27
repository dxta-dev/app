package template

import (
	"bytes"
	"context"
	"io"
	"log"
	"math"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/dxta-dev/app/internal/util"
	"github.com/wcharczuk/go-chart/v2"
)

type TimeSeries struct {
	Title   string
	XValues []float64
	YValues []float64
	Weeks   []string
}

type StartEndWeek struct {
	Start string
	End   string
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

	highest = highest * 1.1

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

	if highest < 4 {
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

type label struct {
	x, y int
	text string
}

func YearLabel(c *chart.Chart, l label, userDefaults ...chart.Style) chart.Renderable {
	return func(r chart.Renderer, box chart.Box, defaults chart.Style) {

		f := util.GetMonospaceFont()

		chart.Draw.Text(r, l.text, l.x, l.y, chart.Style{
			FontColor:           chart.ColorRed,
			FontSize:            12,
			Font:                f,
			TextRotationDegrees: -90,
		})
	}
}

func MonthLabel(c *chart.Chart, l label, userDefaults ...chart.Style) chart.Renderable {
	return func(r chart.Renderer, box chart.Box, defaults chart.Style) {

		f := util.GetMonospaceFont()

		chart.Draw.Text(r, l.text, l.x, l.y, chart.Style{
			FontColor: chart.ColorBlack,
			FontSize:  12,
			Font:      f,
		})
	}
}

func CutOffLabel(c *chart.Chart, l label, userDefaults ...chart.Style) chart.Renderable {
	return func(r chart.Renderer, box chart.Box, defaults chart.Style) {
		f := util.GetMonospaceFont()

		chart.Draw.Text(r, l.text, l.x, l.y, chart.Style{
			FontColor: chart.ColorBlue,
			FontSize:  12,
			Font:      f,
		})
	}
}

func TimeSeriesChart(series TimeSeries, cutoff CutOffWeeks) templ.Component {
	YAxisValues := getYAxisValues(series.YValues)

	f := util.GetMonospaceFont()

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
		Style: chart.Style{
			StrokeWidth: 3,
			StrokeColor: chart.ColorBlue,
		},
	}

	graph := chart.Chart{
		Font: f,
		Background: chart.Style{
			Padding: chart.Box{
				Top:    20,
				Left:   0,
				Right:  0,
				Bottom: 0,
			},
		},
		YAxis: chart.YAxis{
			Ticks: []chart.Tick{
				{Value: YAxisValues[0], Label: util.FormatYAxisValues(YAxisValues[0])},
				{Value: YAxisValues[1], Label: util.FormatYAxisValues(YAxisValues[1])},
				{Value: YAxisValues[2], Label: util.FormatYAxisValues(YAxisValues[2])},
				{Value: YAxisValues[3], Label: util.FormatYAxisValues(YAxisValues[3])},
				{Value: YAxisValues[4], Label: util.FormatYAxisValues(YAxisValues[4])},
			},
		},
		Height: 320,
		Width:  650,
	}

	for i, week := range series.Weeks {
		week_part := strings.Split(week, "-")
		graph.XAxis.Ticks = append(graph.XAxis.Ticks, chart.Tick{
			Value: float64(i),
			Label: week_part[1],
		})
	}

	firstDay, _, err := util.ParseYearWeek(series.Weeks[0])
	if err != nil {
		log.Fatal(err)
	}
	months := util.GetStartOfMonths(series.Weeks)
	monthLabels := []label{}

	for _, startOfMonth := range months {
		xvalue := startOfMonth.Sub(firstDay).Hours() / 24 / 7
		x := int(620 / 12 * xvalue)

		if x < 0 {
			x = 0
		}

		if x > 620 {
			x = 620
		}

		monthLabels = append(monthLabels, label{
			x:    x,
			y:    22,
			text: startOfMonth.Format("Jan"),
		})

		if xvalue <= 0 || xvalue >= float64(len(series.Weeks)) {
			continue
		}

		if startOfMonth.Month() == time.January {
			yearLabel := label{
				x:    x + 38,
				y:    80,
				text: startOfMonth.Format("2006"),
			}
			graph.Elements = append(graph.Elements, YearLabel(&graph, yearLabel))

			prevYearLabel := label{
				x:    x + 16,
				y:    80,
				text: (startOfMonth.AddDate(-1, 0, 0)).Format("2006"),
			}
			graph.Elements = append(graph.Elements, YearLabel(&graph, prevYearLabel))
		}

		if startOfMonth.Month() == time.January {
			gridLine := chart.ContinuousSeries{
				XValues: []float64{xvalue, xvalue},
				YValues: []float64{0, YAxisValues[len(YAxisValues)-1] * 268 / 273},
				Style: chart.Style{
					StrokeWidth: 1.0,
					StrokeColor: chart.ColorRed,
				},
			}
			graph.Series = append(graph.Series, gridLine)
		} else {
			gridLine := chart.ContinuousSeries{
				XValues: []float64{xvalue, xvalue},
				YValues: []float64{0, YAxisValues[len(YAxisValues)-1]},
				Style: chart.Style{
					StrokeWidth:     1.0,
					StrokeColor:     chart.ColorBlack,
					StrokeDashArray: []float64{5, 7},
				},
			}
			graph.Series = append(graph.Series, gridLine)
		}

	}

	for i, monthLabel := range monthLabels {
		if i == len(monthLabels)-1 {
			monthLabel.x = (620-monthLabel.x)/2 + monthLabel.x
		} else {
			monthLabel.x = (monthLabels[i+1].x-monthLabel.x)/2 + monthLabel.x
		}

		graph.Elements = append(graph.Elements, MonthLabel(&graph, monthLabel))
	}

	graph.Series = append(graph.Series, lineSeries)
	graph.Series = append(graph.Series, mainSeries)

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.SVG, buffer)

	if err != nil {
		log.Fatal(err)
	}

	html := buffer.String()

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, html)
		return err
	})
}

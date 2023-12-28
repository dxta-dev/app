package handlers

import (
	"context"
	"dxta-dev/app/internals/templates"
	"fmt"
	"math"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func findYValues(x, y, r, x1 float64) (float64, float64, bool) {
	d := math.Abs(x - x1)
	if d >= r {
		return 0, 0, true
	}
	verticalDistance := math.Sqrt(r*r - d*d)
	y1 := y + verticalDistance
	y2 := y - verticalDistance
	return y1, y2, false
}

func Distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

func Intersect(x1, y1, x2, y2, r float64) bool {
	d := Distance(x1, y1, x2, y2)

	if x1 == x2 && y1 == y2 {
		return true
	}

	if d > 2*r {
		return false
	}
	if d == 2*r {
		return false
	}
	return true
}

func series2() templates.ScatterSeries {
	var xvalues []float64
	var yvalues []float64

	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1).Truncate(24 * time.Hour)

	times := []time.Time{
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 14, 0, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 14, 0, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 14, 0, 0, 0, time.UTC),
		time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 14, 0, 0, 0, time.UTC),
	}

	for _, time := range times {
		xSecondsValue := float64(time.Unix() - startOfWeek.Unix())
		xvalues = append(xvalues, xSecondsValue)
		yvalues = append(yvalues, 60*60*12)
	}

	direction := true

	for i := 0; i < len(xvalues); i++ {
		for j := i + 1; j < len(xvalues); j++ {
			if i == j {
				continue
			}
			x := xvalues[i]
			y := yvalues[i]
			x1 := xvalues[j]
			y1 := yvalues[j]
			r := 2160 * 3.0

			if !Intersect(x, y, x1, y1, r/2) {
				continue
			}

			y1, y2, notFound := findYValues(x, y, r, x1)
			if notFound {
				continue
			}
			if direction {
				yvalues[j] = y1
			} else {
				yvalues[j] = y2
			}
			fmt.Println(i, j, "y1", y1, "y2", y2, direction)

			direction = !direction

		}
	}

	return templates.ScatterSeries{
		Title:   "series 1",
		XValues: xvalues,
		YValues: yvalues,
	}
}

func (a *App) ScatterPacking(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	page := &templates.Page{
		Title:   "Charts",
		Boosted: h.HxBoosted,
	}

	var chartData []templates.ScatterSeries

	chartData = append(chartData, series2())

	components := templates.Scatter(page, chartData)
	return components.Render(context.Background(), c.Response().Writer)
}

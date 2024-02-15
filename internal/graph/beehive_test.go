package graph

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/dxta-dev/app/internal/data"
)

func BenchmarkGenerateHexagonGrid(b *testing.B) {
	height, width, hexHeight, hexWidth, r, rows, cols := setup(1400, 200, 5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateHexagonGrid(width, height, hexWidth, hexHeight, r, rows, cols)
	}
}

func BenchmarkFindNearestHex(b *testing.B) {
	height, width, hexHeight, hexWidth, r, rows, cols := setup(1400, 200, 5)
	hexagons := generateHexagonGrid(width, height, hexWidth, hexHeight, r, rows, cols)
	takenHexagons := make(map[Hexagon]bool)
	x, y := 216000.0, 43200.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findNearestHex(hexagons, takenHexagons, x, y, r)
	}
}

func BenchmarkFindNearestHexTaken20Percent(b *testing.B) {
	height, width, hexHeight, hexWidth, r, rows, cols := setup(1400, 200, 5)
	hexagons := generateHexagonGrid(width, height, hexWidth, hexHeight, r, rows, cols)
	takenHexagons := make(map[Hexagon]bool)
	x, y := 216000.0, 43200.0

	for i, hex := range hexagons {
		if i%5 == 0 {
			takenHexagons[hex] = true
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findNearestHex(hexagons, takenHexagons, x, y, r)
	}
}

func BenchmarkFindNearestHexTaken33Percent(b *testing.B) {
	height, width, hexHeight, hexWidth, r, rows, cols := setup(1400, 200, 5)
	hexagons := generateHexagonGrid(width, height, hexWidth, hexHeight, r, rows, cols)
	takenHexagons := make(map[Hexagon]bool)
	x, y := 216000.0, 43200.0

	for i, hex := range hexagons {
		if i%3 == 0 {
			takenHexagons[hex] = true
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findNearestHex(hexagons, takenHexagons, x, y, r)
	}
}

func BenchmarkFindNearestHexTaken(b *testing.B) {
	height, width, hexHeight, hexWidth, r, rows, cols := setup(1400, 200, 5)
	hexagons := generateHexagonGrid(width, height, hexWidth, hexHeight, r, rows, cols)
	takenHexagons := make(map[Hexagon]bool)
	x, y := 216000.0, 43200.0

	for i, hex := range hexagons {
		if i == 0 {
			continue
		}
		takenHexagons[hex] = true
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findNearestHex(hexagons, takenHexagons, x, y, r)
	}
}

func BenchmarkBeehive(b *testing.B) {
	var times []time.Time

	var xvalues []float64
	var yvalues []float64

	sort.Sort(data.DataList)

	startOfWeek := time.Unix(1696204800, 0)

	for _, d := range data.DataList {
		t := time.Unix(d.Timestamp/1000, 0)
		times = append(times, t)
	}

	for _, t := range times {
		xSecondsValue := float64(t.Unix() - startOfWeek.Unix())
		xvalues = append(xvalues, xSecondsValue)
		yvalues = append(yvalues, 60*60*12)
	}

	chartWidth, chartHeight, dotWidth := 1400, 200, 5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Beehive(xvalues, yvalues, chartWidth, chartHeight, dotWidth)
	}
}

func TestBeehiveFunctionality(t *testing.T) {
	tests := []struct {
		name             string
		xValues, yValues []float64
		chartWidth       int
		chartHeight      int
		dotWidth         int
		expectedXValues  []float64
		expectedYValues  []float64
	}{
		{
			name:            "SimpleCase",
			xValues:         []float64{100, 200},
			yValues:         []float64{100, 200},
			chartWidth:      400,
			chartHeight:     400,
			dotWidth:        10,
			expectedXValues: []float64{0, 8978.95138643706},
			expectedYValues: []float64{0, 5184},
		},
		{
			name:            "OverlapCase",
			xValues:         []float64{100, 105},
			yValues:         []float64{100, 105},
			chartWidth:      400,
			chartHeight:     400,
			dotWidth:        10,
			expectedXValues: []float64{0, 8978.95138643706},
			expectedYValues: []float64{0, 5184},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotXValues, gotYValues := Beehive(tt.xValues, tt.yValues, tt.chartWidth, tt.chartHeight, tt.dotWidth)

			if !reflect.DeepEqual(gotXValues, tt.expectedXValues) || !reflect.DeepEqual(gotYValues, tt.expectedYValues) {
				t.Errorf("Beehive() = (%v, %v), want (%v, %v)", gotXValues, gotYValues, tt.expectedXValues, tt.expectedYValues)
			}
		})
	}
}

func TestBeehiveEdgeCases(t *testing.T) {
	tests := []struct {
		name             string
		xValues, yValues []float64
		chartWidth       int
		chartHeight      int
		dotWidth         int
		expectedXValues  []float64
		expectedYValues  []float64
	}{
		{
			name:            "EmptyArrays",
			xValues:         []float64{},
			yValues:         []float64{},
			chartWidth:      400,
			chartHeight:     400,
			dotWidth:        10,
			expectedXValues: []float64{},
			expectedYValues: []float64{},
		},
		{
			name:            "SmallDotWidth",
			xValues:         []float64{100, 200},
			yValues:         []float64{100, 200},
			chartWidth:      400,
			chartHeight:     400,
			dotWidth:        1,
			expectedXValues: []float64{0, 897.8951386437059},
			expectedYValues: []float64{0, 518.4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotXValues, gotYValues := Beehive(tt.xValues, tt.yValues, tt.chartWidth, tt.chartHeight, tt.dotWidth)

			if !reflect.DeepEqual(gotXValues, tt.expectedXValues) || !reflect.DeepEqual(gotYValues, tt.expectedYValues) {
				t.Errorf("Beehive() = (%v, %v), want (%v, %v)", gotXValues, gotYValues, tt.expectedXValues, tt.expectedYValues)
			}
		})
	}
}

func TestBeehiveConsistency(t *testing.T) {
	xValues := []float64{100, 200, 300}
	yValues := []float64{100, 200, 300}
	chartWidth := 500
	chartHeight := 500
	dotWidth := 10

	firstRunX, firstRunY := Beehive(xValues, yValues, chartWidth, chartHeight, dotWidth)
	secondRunX, secondRunY := Beehive(xValues, yValues, chartWidth, chartHeight, dotWidth)
	thirdRunX, thirdRunY := Beehive(xValues, yValues, chartWidth, chartHeight, dotWidth)

	if !reflect.DeepEqual(firstRunX, secondRunX) || !reflect.DeepEqual(firstRunY, secondRunY) ||
		!reflect.DeepEqual(firstRunX, thirdRunX) || !reflect.DeepEqual(firstRunY, thirdRunY) {
		t.Error("Beehive() produced inconsistent results over multiple runs")
	}
}

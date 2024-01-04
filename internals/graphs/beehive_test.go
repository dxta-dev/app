package graphs

import (
	"math"
	"reflect"
	"testing"
	"time"
	"sort"

	"dxta-dev/app/internals/data"
)

func BenchmarkGenerateHexagonGrid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateHexagonGrid(1400, 200, 5)
	}
}

func BenchmarkFindNearestHex(b *testing.B) {
	hexagons := generateHexagonGrid(1400, 200, 5)
	takenHexagons := make(map[Hexagon]bool)
	x, y := 216000.0, 43200.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findNearestHex(hexagons, takenHexagons, x, y)
	}
}

func BenchmarkFindNearestHexTaken(b *testing.B) {
	hexagons := generateHexagonGrid(1400, 200, 5)
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
		findNearestHex(hexagons, takenHexagons, x, y)
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

func TestGenerateHexagonGrid(t *testing.T) {
	tests := []struct {
		name           string
		width, height  float64
		r              float64
		expectedLength int
	}{
		{"SmallGrid", 1000, 1000, 50, 120},
		{"MediumGrid", 2000, 1500, 75, 160},
		{"LargeGrid", 3000, 2500, 100, 225},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hexagons := generateHexagonGrid(tt.width, tt.height, tt.r)
			if len(hexagons) != tt.expectedLength {
				t.Errorf("generateHexagonGrid(%v, %v, %v) got %v hexagons, want %v", tt.width, tt.height, tt.r, len(hexagons), tt.expectedLength)
			}

			for _, hex := range hexagons {
				if hex.X < 0 || hex.X > tt.width || hex.Y < 0 || hex.Y > tt.height {
					t.Errorf("Hexagon %v is out of bounds in grid size %v x %v", hex, tt.width, tt.height)
				}
			}
		})
	}
}

func TestGenerateHexagonGridBoundary(t *testing.T) {
	tests := []struct {
		name           string
		width, height  float64
		r              float64
		expectedLength int
	}{
		{"VerySmallGrid", 10, 10, 5, 2},
		{"ZeroGrid", 0, 0, 50, 0},
		{"NegativeDimensions", -100, -100, 50, 0},
		{"LargeGridSmallRadius", 10000, 10000, 1, 28870000},
		{"SmallGridLargeRadius", 500, 500, 250, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hexagons := generateHexagonGrid(tt.width, tt.height, tt.r)
			if len(hexagons) != tt.expectedLength {
				t.Errorf("generateHexagonGrid(%v, %v, %v) got %v hexagons, want %v", tt.width, tt.height, tt.r, len(hexagons), tt.expectedLength)
			}

			for _, hex := range hexagons {
				if hex.X < 0 || hex.X > tt.width || hex.Y < 0 || hex.Y > tt.height {
					t.Errorf("Hexagon %v is out of bounds in grid size %v x %v", hex, tt.width, tt.height)
				}
			}
		})
	}
}

func TestRemoveTakenHexagons(t *testing.T) {
	hexagons := []Hexagon{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 2},
		{X: 3, Y: 3},
	}

	takenHexagons := map[Hexagon]bool{
		{X: 1, Y: 1}: true,
		{X: 3, Y: 3}: true,
	}

	result := removeTakenHexagons(hexagons, takenHexagons)

	expected := []Hexagon{
		{X: 0, Y: 0},
		{X: 2, Y: 2},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("removeTakenHexagons() = %v, want %v", result, expected)
	}
}

func TestRemoveTakenHexagonsEmpty(t *testing.T) {
	hexagons := []Hexagon{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 2},
	}

	takenHexagons := make(map[Hexagon]bool)

	result := removeTakenHexagons(hexagons, takenHexagons)

	if !reflect.DeepEqual(result, hexagons) {
		t.Errorf("removeTakenHexagons() with empty map = %v, want %v", result, hexagons)
	}
}

func TestRemoveTakenHexagonsEmptyInverse(t *testing.T) {
	hexagons := []Hexagon{}

	takenHexagons := map[Hexagon]bool{
		{X: 1, Y: 1}: true,
		{X: 3, Y: 3}: true,
	}

	result := removeTakenHexagons(hexagons, takenHexagons)

	if !reflect.DeepEqual(result, hexagons) {
		t.Errorf("removeTakenHexagons() with empty hexagons = %v, want %v", result, hexagons)
	}
}

func TestDistanceCalculation(t *testing.T) {
	tests := []struct {
		name     string
		hexagon  Hexagon
		x, y     float64
		expected float64
	}{
		{"ZeroDistance", Hexagon{X: 0, Y: 0}, 0, 0, 0},
		{"PositiveCoordinates", Hexagon{X: 3, Y: 4}, 0, 0, 5},
		{"NegativeCoordinates", Hexagon{X: -3, Y: -4}, 0, 0, 5},
		{"MixedCoordinates", Hexagon{X: -3, Y: 4}, 0, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := distance(tt.hexagon, tt.x, tt.y)
			if math.Abs(result-tt.expected) > 1e-6 {
				t.Errorf("distance(%v, %v, %v) = %v, want %v", tt.hexagon, tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

func TestFindNearestHex(t *testing.T) {
	hexagons := []Hexagon{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 2},
		{X: 3, Y: 3},
	}

	takenHexagons := map[Hexagon]bool{
		{X: 0, Y: 0}: true,
		{X: 1, Y: 1}: true,
	}

	tests := []struct {
		name     string
		x, y     float64
		expected Hexagon
	}{
		{"ClosestToOrigin", 0, 0, Hexagon{X: 2, Y: 2}},
		{"ClosestToMidPoint", 1.5, 1.5, Hexagon{X: 2, Y: 2}},
		{"ClosestToFarPoint", 3, 3, Hexagon{X: 3, Y: 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findNearestHex(hexagons, takenHexagons, tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("findNearestHex() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFindNearestHexAllTaken(t *testing.T) {
	hexagons := []Hexagon{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 2},
	}

	takenHexagons := map[Hexagon]bool{
		{X: 0, Y: 0}: true,
		{X: 1, Y: 1}: true,
		{X: 2, Y: 2}: true,
	}

	x, y := 0.5, 0.5

	result := findNearestHex(hexagons, takenHexagons, x, y)

	expected := Hexagon{}

	if result != expected {
		t.Errorf("findNearestHex() with all hexagons taken = %v, want %v", result, expected)
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

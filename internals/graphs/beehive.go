package graphs

import (
	"math"
)

var unit float64 = 432.0

type Hexagon struct {
	X, Y float64
}

func generateHexagonGrid(width, height, r float64) []Hexagon {
	var hexagons []Hexagon

	hexHeight := 2 * r
	hexWidth := 4 * math.Sqrt(3) * r / 3

	offset := hexHeight / 2

	rows := int(height/hexHeight) + 1
	cols := int(width/hexWidth*4/3) + 1

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			x := float64(col) * hexWidth * 3 / 4
			y := float64(row) * hexHeight
			if col%2 == 1 {
				y += offset
			}

			if x < width && y < height {
				hexagons = append(hexagons, Hexagon{x, y})
			}

		}
	}
	return hexagons

}

func removeTakenHexagons(hexagons []Hexagon, takenHexagons map[Hexagon]bool) []Hexagon {
	result := []Hexagon{}

	for _, hexagon := range hexagons {
		if !takenHexagons[hexagon] {
			result = append(result, hexagon)
		}
	}

	return result
}

func distance(h Hexagon, x, y float64) float64 {
	return math.Sqrt(math.Pow(h.X-x, 2) + math.Pow(h.Y-y, 2))
}

func findNearestHex(hexagons []Hexagon, takenHexagons map[Hexagon]bool, x, y, r float64) Hexagon {
	minDistanceSquared := math.MaxFloat64
	var nearestHex Hexagon

	for _, hex := range hexagons {
		if takenHexagons[hex] {
			continue
		}

		dX, dY := hex.X-x, hex.Y-y
		distanceSquared := dX*dX + dY*dY

		if distanceSquared < minDistanceSquared {
			minDistanceSquared = distanceSquared
			nearestHex = hex

			if minDistanceSquared < (r / 2) * (r / 2) {
				break
			}
		}
	}

	return nearestHex
}

func Beehive(xValues []float64, yValues []float64, chartWidth, chartHeight, dotWidth int) ([]float64, []float64) {
	radius := float64(dotWidth) * unit
	r := radius * 1.2

	hexagons := generateHexagonGrid(float64(chartWidth)*unit, float64(chartHeight)*unit, r)

	takenHex := make(map[Hexagon]bool)

	for i := 0; i < len(xValues); i++ {
		x := xValues[i]
		y := yValues[i]

		for _, hexagon := range hexagons {
			if math.Abs(x-hexagon.X) < r && math.Abs(y-hexagon.Y) < r {
				if _, exists := takenHex[hexagon]; exists {
					nearHex := findNearestHex(hexagons, takenHex, x, y, r)
					xValues[i] = nearHex.X
					yValues[i] = nearHex.Y
					takenHex[nearHex] = true
					break
				} else {
					xValues[i] = hexagon.X
					yValues[i] = hexagon.Y
					takenHex[hexagon] = true
					break
				}
			}
		}
	}

	xvalues := []float64{}
	yvalues := []float64{}

	for i := 0; i < len(xValues); i++ {
		xvalues = append(xvalues, xValues[i])
		yvalues = append(yvalues, yValues[i])
	}

	return xvalues, yvalues
}

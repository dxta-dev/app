package util

import (
	"fmt"
	"math"
	"strings"
)

func FormatYAxisValues(num float64) string {
	absNum := math.Abs(num)
	var formatted string

	switch {
	case absNum < 1000:
		if num == math.Floor(num) {
			formatted = fmt.Sprintf("%d", int(num))
		} else {
			formatted = fmt.Sprintf("%.1f", num)
		}
	case absNum >= 1000 && absNum < 1000000:
		formatted = fmt.Sprintf("%.1fK", num/1000)
	case absNum >= 1000000:
		formatted = fmt.Sprintf("%.1fM", num/1000000)
	default:
		formatted = fmt.Sprintf("%f", num)
	}

	if len(formatted) < 5 {
		formatted += strings.Repeat("\u00a0", 5-len(formatted))
	}

	return formatted
}

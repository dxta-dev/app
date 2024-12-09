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
	case absNum >= 1000 && absNum < 100000:
		formatted = fmt.Sprintf("%.1fK", num/1000)
	case absNum >= 100000 && absNum < 1000000:
		rounded := math.Round(num / 1000)
		formatted = fmt.Sprintf("%.0fK", rounded)
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

func GetTimeAxisNormalizationFactorLabel(factor float64) string {
	const minute = 60 * 1000
	const hour = 60 * minute
	const day = 24 * hour
	const week = 7 * day

	if factor == 1.0/minute {
		return "minutes"
	}

	if factor == 1.0/hour {
		return "hours"
	}

	if factor == 1.0/day {
		return "days"
	}

	if factor == 1.0/week {
		return "weeks"
	}

	return "unknown"
}

func GetTimeAxisNormalizationFactor(milis float64) float64 {
	const minute = 60 * 1000
	const hour = 60 * minute
	const day = 24 * hour
	const week = 7 * day
	const twoWeeks = 2 * week

	if milis < hour {
		return 1.0 / minute
	}

	if milis < day {
		return 1.0 / hour
	}

	if milis < twoWeeks {
		return 1.0 / day
	}

	return 1.0 / week
}

func FormatTimeAxisValue(milis float64) string {
	if milis == 0.0 {
		return "0"
	}

	const minute = 60 * 1000
	const hour = 60 * minute

	factor := GetTimeAxisNormalizationFactor(milis)
	factorLabel := GetTimeAxisNormalizationFactorLabel(factor)

	if milis < hour {
		return "less than an hour"
	}

	return fmt.Sprintf("%.1f %s", milis*factor, factorLabel)
}

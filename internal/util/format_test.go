package util

import (
	"testing"
	"strings"
)

func TestFormatYAxisValues(t *testing.T) {
	nbsp := "\u00a0"
	testCases := []struct {
		name     string
		num      float64
		expected string
	}{
		{"Zero", 0, "0" + strings.Repeat(nbsp, 4)},
		{"Less than 10 with no decimals", 5, "5" + strings.Repeat(nbsp, 4)},
		{"Less than 10 with decimals", 7.2, "7.2" + strings.Repeat(nbsp, 2)},
		{"Less than 100 with no decimals", 50, "50" + strings.Repeat(nbsp, 3)},
		{"Less than 1000 with no decimals", 500, "500" + strings.Repeat(nbsp, 2)},
		{"Exactly 1000", 1000, "1.0K" + nbsp},
		{"Less than 10000 with no decimals", 5000, "5.0K" + nbsp},
		{"Less than 100000 with no decimals", 50000, "50.0K"},
		{"Exactly 1 Million", 1000000, "1.0M" + nbsp},
		{"MIllions with no decimals", 2000000, "2.0M" + nbsp},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatYAxisValues(tc.num)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

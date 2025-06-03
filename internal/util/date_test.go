package util

import (
	"reflect"
	"testing"
	"time"
)

func TestGetFormattedWeek(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		input    time.Time
	}{
		{
			name:     "standard date",
			expected: "2024-W20",
			input:    time.Date(2024, time.May, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "beginning of the year",
			expected: "2024-W01",
			input:    time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "end of the year",
			expected: "2025-W01",
			input:    time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFormattedWeek(tt.input)
			if tt.expected != got {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestGetLastNWeeks(t *testing.T) {
	t.Run("current date", func(t *testing.T) {
		n := 4
		weeks := GetLastNWeeks(time.Now(), n)
		if len(weeks) != n {
			t.Errorf("Expected %d weeks, got %d", n, len(weeks))
		}
	})

	t.Run("specific date", func(t *testing.T) {
		date := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC) // A known date
		n := 2
		expectedWeeks := []string{"2022-W52", "2023-W01"}
		weeks := GetLastNWeeks(date, n)
		for i, week := range weeks {
			if week != expectedWeeks[i] {
				t.Errorf("Expected %s, got %s", expectedWeeks[i], week)
			}
		}
	})

	t.Run("n is zero", func(t *testing.T) {
		weeks := GetLastNWeeks(time.Now(), 0)
		if len(weeks) != 0 {
			t.Errorf("Expected 0 weeks, got %d", len(weeks))
		}
	})
}

func TestParseISOWeek(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ISOWeek
	}{
		{
			name:     "valid ISO week",
			input:    "2024-W25",
			expected: ISOWeek{Year: 2024, Week: 25},
		},
		{
			name:     "out of range ISO week",
			input:    "2024-W54",
			expected: ISOWeek{},
		},
		{
			name:     "invalid input",
			input:    "string",
			expected: ISOWeek{},
		},
		{
			name:     "more than 2 parts parsed",
			input:    "2024-W01-W01",
			expected: ISOWeek{},
		},
		{
			name:     "wrong format of input",
			input:    "-W-2024",
			expected: ISOWeek{},
		},
		{
			name:     "valid leap week",
			input:    "2020-W53",
			expected: ISOWeek{Year: 2020, Week: 53},
		},
		{
			name:     "invalid leap week",
			input:    "2019-W53",
			expected: ISOWeek{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ParseISOWeek(tt.input)
			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestSortISOWeeks(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "sort by years",
			input:    []string{"2024-W33", "2021-W33", "2025-W33"},
			expected: []string{"2021-W33", "2024-W33", "2025-W33"},
		},
		{
			name:     "sort by weeks",
			input:    []string{"2024-W33", "2024-W01", "2024-W20"},
			expected: []string{"2024-W01", "2024-W20", "2024-W33"},
		},
		{
			name: "sort by years and weeks",
			input: []string{
				"2024-W33",
				"2024-W01",
				"2024-W20",
				"2024-W33",
				"2021-W33",
				"2025-W33",
			},
			expected: []string{
				"2021-W33",
				"2024-W01",
				"2024-W20",
				"2024-W33",
				"2024-W33",
				"2025-W33",
			},
		},
		{
			name:     "invalid input",
			input:    []string{"2024W33", "2024-W01", "200", "2024-W33", "2021-W33", "2025-W33"},
			expected: []string{"2021-W33", "2024-W01", "2024-W33", "2025-W33"},
		},
		{
			name: "valid leap weeks",
			input: []string{
				"2020-W53",
				"2021-W53",
				"2022-W53",
				"2023-W53",
				"2024-W53",
				"2025-W53",
				"2026-W53",
				"2027-W53",
				"2028-W53",
				"2029-W53",
				"2030-W53",
				"2031-W53",
				"2032-W53",
			},
			expected: []string{"2020-W53", "2026-W53", "2032-W53"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortISOWeeks(tt.input)
			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

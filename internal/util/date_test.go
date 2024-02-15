package util

import (
	"testing"
	"time"
)

func TestGetStartOfTheWeek(t *testing.T) {
	tests := []struct {
		name     string
		expected time.Time
		input    time.Time
	}{
		{
			name:     "middle of the week",
			expected: time.Date(2024, time.January, 8, 0, 0, 0, 0, time.UTC),
			input:    GetStartOfTheWeek(time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:     "beginning of the week",
			expected: time.Date(2024, time.January, 8, 0, 0, 0, 0, time.UTC),
			input:    GetStartOfTheWeek(time.Date(2024, time.January, 8, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:     "end of the week",
			expected: time.Date(2024, time.January, 8, 0, 0, 0, 0, time.UTC),
			input:    GetStartOfTheWeek(time.Date(2024, time.January, 14, 0, 0, 0, 0, time.UTC)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStartOfTheWeek(tt.input)
			if tt.expected != got {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

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

func TestParseYearWeek(t *testing.T) {
	tests := []struct {
		name     string
		expected time.Time
		input    string
	}{
		{
			name: "first week of 1980",
			expected: time.Date(1979, time.December, 31, 0, 0, 0, 0, time.UTC),
			input: "1980-W01",
		},
		{
			name: "first week of 1981",
			expected: time.Date(1980, time.December, 29, 0, 0, 0, 0, time.UTC),
			input: "1981-W01",
		},
		{
			name: "first week of 2015",
			expected: time.Date(2014, time.December, 29, 0, 0, 0, 0, time.UTC),
			input: "2015-W01",
		},
		{
			name:     "first week of 2022",
			expected: time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
			input:    "2022-W01",
		},
		{
			name:     "last week of 2022",
			expected: time.Date(2022, time.December, 26, 0, 0, 0, 0, time.UTC),
			input:    "2022-W52",
		},
		{
			name:     "sunday?",
			expected: time.Date(2023, time.December, 18, 0, 0, 0, 0, time.UTC),
			input:    "2023-W51",
		},
		{
			name: "begging of the last year",
			expected: time.Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC),
			input: "2023-W01",
		},
		{
			name:     "beginning of the year",
			expected: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			input:    "2024-W01",
		},
		{
			name:     "valid year and week",
			expected: time.Date(2024, time.May, 13, 0, 0, 0, 0, time.UTC),
			input:    "2024-W20",
		},
		{
			name:     "invalid format",
			expected: time.Time{},
			input:    "2024W20",
		},
		{
			name:     "non existent week number",
			expected: time.Time{},
			input:    "2024-W54",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ParseYearWeek(tt.input)
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

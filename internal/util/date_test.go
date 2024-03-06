package util

import (
	"reflect"
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

func TestGetStartOfTheMonth(t *testing.T) {
	tests := []struct {
		name     string
		expected []time.Time
		input    []string
	}{
		{
			name: "Standard last 12 weeks",
			expected: []time.Time{
				parseTime("2023-12-01T00:00:00Z"),
				parseTime("2024-01-01T00:00:00Z"),
				parseTime("2024-02-01T00:00:00Z"),
			},
			input: []string{
				"2023-W50", "2023-W51", "2023-W52",
				"2024-W01", "2024-W02", "2024-W03",
				"2024-W04", "2024-W05", "2024-W06",
				"2024-W07", "2024-W08", "2024-W09",
			},
		},
		{
			name: "Last 96 weeks",
			expected: []time.Time{
				parseTime("2022-05-01T00:00:00Z"), parseTime("2022-06-01T00:00:00Z"), parseTime("2022-07-01T00:00:00Z"), parseTime("2022-08-01T00:00:00Z"),
				parseTime("2022-09-01T00:00:00Z"), parseTime("2022-10-01T00:00:00Z"), parseTime("2022-11-01T00:00:00Z"), parseTime("2022-12-01T00:00:00Z"),
				parseTime("2023-01-01T00:00:00Z"), parseTime("2023-02-01T00:00:00Z"), parseTime("2023-03-01T00:00:00Z"), parseTime("2023-04-01T00:00:00Z"),
				parseTime("2023-05-01T00:00:00Z"), parseTime("2023-06-01T00:00:00Z"), parseTime("2023-07-01T00:00:00Z"), parseTime("2023-08-01T00:00:00Z"),
				parseTime("2023-09-01T00:00:00Z"), parseTime("2023-10-01T00:00:00Z"), parseTime("2023-11-01T00:00:00Z"), parseTime("2023-12-01T00:00:00Z"),
				parseTime("2024-01-01T00:00:00Z"), parseTime("2024-02-01T00:00:00Z"),
			},
			input: []string{
				"2022-W18", "2022-W19", "2022-W20", "2022-W21", "2022-W22", "2022-W23", "2022-W24", "2022-W25", "2022-W26", "2022-W27",
				"2022-W28", "2022-W29", "2022-W30", "2022-W31", "2022-W32", "2022-W33", "2022-W34", "2022-W35", "2022-W36", "2022-W37",
				"2022-W38", "2022-W39", "2022-W40", "2022-W41", "2022-W42", "2022-W43", "2022-W44", "2022-W45", "2022-W46", "2022-W47",
				"2022-W48", "2022-W49", "2022-W50", "2022-W51", "2022-W52", "2023-W01", "2023-W02", "2023-W03", "2023-W04", "2023-W05",
				"2023-W06", "2023-W07", "2023-W08", "2023-W09", "2023-W10", "2023-W11", "2023-W12", "2023-W13", "2023-W14", "2023-W15",
				"2023-W16", "2023-W17", "2023-W18", "2023-W19", "2023-W20", "2023-W21", "2023-W22", "2023-W23", "2023-W24", "2023-W25",
				"2023-W26", "2023-W27", "2023-W28", "2023-W29", "2023-W30", "2023-W31", "2023-W32", "2023-W33", "2023-W34", "2023-W35",
				"2023-W36", "2023-W37", "2023-W38", "2023-W39", "2023-W40", "2023-W41", "2023-W42", "2023-W43", "2023-W44", "2023-W45",
				"2023-W46", "2023-W47", "2023-W48", "2023-W49", "2023-W50", "2023-W51", "2023-W52", "2024-W01", "2024-W02", "2024-W03",
				"2024-W04", "2024-W05", "2024-W06", "2024-W07", "2024-W08", "2024-W09",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStartOfMonths(tt.input)
			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}

}

func parseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		panic(err)
	}
	return t
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
			name:     "first week of 1980",
			expected: time.Date(1979, time.December, 31, 0, 0, 0, 0, time.UTC),
			input:    "1980-W01",
		},
		{
			name:     "first week of 1981",
			expected: time.Date(1980, time.December, 29, 0, 0, 0, 0, time.UTC),
			input:    "1981-W01",
		},
		{
			name:     "first week of 2015",
			expected: time.Date(2014, time.December, 29, 0, 0, 0, 0, time.UTC),
			input:    "2015-W01",
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
			name:     "begging of the last year",
			expected: time.Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC),
			input:    "2023-W01",
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

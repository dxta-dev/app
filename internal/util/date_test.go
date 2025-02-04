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
				time.Date(2023, time.December, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
			input: []string{
				"2023-W50", "2023-W51", "2023-W52",
				"2024-W01", "2024-W02", "2024-W03",
				"2024-W04", "2024-W05", "2024-W06",
				"2024-W07", "2024-W08", "2024-W09",
			},
		},
		{
			name: "Overlap of 2022/2023 13 weeks",
			expected: []time.Time{
				time.Date(2022, time.December, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC),
			},
			input: []string{
				"2022-W50", "2022-W51", "2022-W52",
				"2023-W01", "2023-W02", "2023-W03",
				"2023-W04", "2023-W05", "2023-W06",
				"2023-W07", "2023-W08", "2023-W09",
				"2023-W10",
			},
		},
		{
			name: "Overlap of 1980/1981 4 weeks",
			expected: []time.Time{
				time.Date(1980, time.December, 1, 0, 0, 0, 0, time.UTC),
				time.Date(1981, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			input: []string{
				"1980-W49", "1980-W50", "1980-W51",
				"1980-W52", "1981-W01", "1981-W02",
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
		name         string
		expectedFrom time.Time
		expectedTo   time.Time
		input        string
	}{
		{
			name:         "first week of 1980",
			expectedFrom: time.Date(1979, time.December, 31, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(1980, time.January, 6, 0, 0, 0, 0, time.UTC),
			input:        "1980-W01",
		},
		{
			name:         "first week of 1981",
			expectedFrom: time.Date(1980, time.December, 29, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(1981, time.January, 4, 0, 0, 0, 0, time.UTC),
			input:        "1981-W01",
		},
		{
			name:         "first week of 2015",
			expectedFrom: time.Date(2014, time.December, 29, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2015, time.January, 4, 0, 0, 0, 0, time.UTC),
			input:        "2015-W01",
		},
		{
			name:         "first week of 2022",
			expectedFrom: time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2022, time.January, 9, 0, 0, 0, 0, time.UTC),
			input:        "2022-W01",
		},
		{
			name:         "last week of 2022",
			expectedFrom: time.Date(2022, time.December, 26, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			input:        "2022-W52",
		},
		{
			name:         "sunday?",
			expectedFrom: time.Date(2023, time.December, 18, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2023, time.December, 24, 0, 0, 0, 0, time.UTC),
			input:        "2023-W51",
		},
		{
			name:         "begging of the last year",
			expectedFrom: time.Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2023, time.January, 8, 0, 0, 0, 0, time.UTC),
			input:        "2023-W01",
		},
		{
			name:         "beginning of the year",
			expectedFrom: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2024, time.January, 7, 0, 0, 0, 0, time.UTC),
			input:        "2024-W01",
		},
		{
			name:         "valid year and week",
			expectedFrom: time.Date(2024, time.May, 13, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2024, time.May, 19, 0, 0, 0, 0, time.UTC),
			input:        "2024-W20",
		},
		{
			name:         "invalid format",
			expectedFrom: time.Time{},
			input:        "2024W20",
		},
		{
			name:         "non existent week number",
			expectedFrom: time.Time{},
			input:        "2024-W54",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFrom, gotTo, _ := ParseYearWeek(tt.input)
			if tt.expectedFrom != gotFrom {
				t.Errorf("expected %v, got %v", tt.expectedFrom, gotFrom)
			}
			if tt.expectedTo != gotTo {
				t.Errorf("expected to %v, got %v", tt.expectedTo, gotTo)
			}
		})
	}
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

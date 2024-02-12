package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetLastNWeeks(date time.Time, n int) []string {
	weeks := make([]string, n)

	for i := 0; i < n; i++ {
		weeks[i] = GetFormattedWeek(date.AddDate(0, 0, -7*i))
	}

	return weeks
}

func GetFormattedWeek(date time.Time) string {
	year, week := date.ISOWeek()

	formattedWeek := fmt.Sprintf("%d-W%02d", year, week)

	return formattedWeek
}

func GetPrevNextWeek(date time.Time) (string, string) {
	startOfWeek := GetStartOfTheWeek(date)
	prevStartOfWeek := startOfWeek.AddDate(0, 0, -7)
	nextStartOfWeek := startOfWeek.AddDate(0, 0, 7)

	prevWeek := GetFormattedWeek(prevStartOfWeek)
	nextWeek := GetFormattedWeek(nextStartOfWeek)

	return prevWeek, nextWeek
}

func GetStartOfTheWeek(date time.Time) time.Time {
	offset := int(time.Monday - date.Weekday())

	if offset > 0 {
		offset = -6
	}

	startOfWeek := date.AddDate(0, 0, offset)

	startOfWeek = startOfWeek.Truncate(24 * time.Hour)

	return startOfWeek
}


func ParseYearWeek(yw string) (time.Time, error) {
	parts := strings.Split(yw, "-W")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid format")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, err
	}

	week, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, err
	}

	firstDayOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)

	if firstDayOfYear.Weekday() != time.Monday {
		for i := 1; i < 4; i++ {
			p := firstDayOfYear.AddDate(0, 0, -i)
			n := firstDayOfYear.AddDate(0, 0, i)

			if p.Weekday() == time.Monday {
				firstDayOfYear = p
				break
			}

			if n.Weekday() == time.Monday {
				firstDayOfYear = n
				break
			}
		}
	}

	startOfWeek := firstDayOfYear.AddDate(0, 0, (week-1)*7)

	if startOfWeek.Year() > year {
		return time.Time{}, fmt.Errorf("invalid week")
	}

	return startOfWeek, nil
}

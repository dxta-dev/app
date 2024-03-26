package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetStartOfMonths(weeks []string) []time.Time {
	monthMap := make(map[time.Time]bool)
	startOfMonths := make([]time.Time, 0)

	for _, week := range weeks {
		date, _, err := ParseYearWeek(week)
		if err != nil {
			continue
		}

		startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
		if _, ok := monthMap[startOfMonth]; !ok {
			monthMap[startOfMonth] = true
			startOfMonths = append(startOfMonths, startOfMonth)
		}

	}
	return startOfMonths
}

func GetLastNWeeks(date time.Time, n int) []string {
	weeks := make([]string, n)

	for i := n; i > 0; i-- {
		weeks[n-i] = GetFormattedWeek(date.AddDate(0, 0, -7*i))
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

func ParseYearWeek(yw string) (time.Time, time.Time, error) {
	parts := strings.Split(yw, "-W")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid format")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	week, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
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
	endOfWeek := startOfWeek.AddDate(0, 0, 6) // End of the week is 6 days after start of the week

	if startOfWeek.Year() > year || endOfWeek.Year() < year {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid week")
	}

	return startOfWeek, endOfWeek, nil
}

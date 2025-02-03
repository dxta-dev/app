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

func GetStartEndWeekDates(yearWeek string) (string, error) {
	var year, week int
	_, err := fmt.Sscanf(yearWeek, "%4d-W%2d", &year, &week)
	if err != nil {
		return "", fmt.Errorf("invalid format")
	}

	firstDayOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)

	daysOffset := (week - 1) * 7
	firstDayOfWeek := firstDayOfYear.AddDate(0, 0, daysOffset)

	for firstDayOfWeek.Weekday() != time.Monday {
		firstDayOfWeek = firstDayOfWeek.AddDate(0, 0, -1)
	}

	lastDayOfWeek := firstDayOfWeek.AddDate(0, 0, 6)

	firstDateStr := firstDayOfWeek.Format("Jan 02")
	lastDateStr := lastDayOfWeek.Format("Jan 02")

	return fmt.Sprintf("%s - %s", firstDateStr, lastDateStr), nil
}

func GetWeeksBetween(startWeek, endWeek string) ([]string, error) {
	var weeks []string
	currentWeek := GetFormattedWeek(time.Now())

	if startWeek == "" && endWeek == "" {
		return nil, fmt.Errorf("either startWeek or endWeek must be provided")
	}

	maxWeeks := 12

	if startWeek != "" {
		startDate, _, err := ParseYearWeek(startWeek)
		if err != nil {
			return nil, fmt.Errorf("invalid start week: %w", err)
		}

		endDate := time.Now()
		if endWeek != "" {
			endDate, _, err = ParseYearWeek(endWeek)
			if err != nil {
				return nil, fmt.Errorf("invalid end week: %w", err)
			}
		}

		currentDate := startDate
		count := 0
		for !currentDate.After(endDate) && count < maxWeeks {
			weeks = append(weeks, GetFormattedWeek(currentDate))
			currentDate = currentDate.AddDate(0, 0, 7)
			count++
		}
	} else if endWeek != "" {
		endDate, _, err := ParseYearWeek(endWeek)
		if err != nil {
			return nil, fmt.Errorf("invalid end week: %w", err)
		}

		startDate := endDate.AddDate(0, 0, (-maxWeeks+1)*7)
		currentDate := startDate
		count := 0
		for !currentDate.After(endDate) && GetFormattedWeek(currentDate) <= currentWeek && count < maxWeeks {
			weeks = append(weeks, GetFormattedWeek(currentDate))
			currentDate = currentDate.AddDate(0, 0, 7)
			count++
		}
	}

	return weeks, nil
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

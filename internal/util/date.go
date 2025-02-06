package util

import (
	"fmt"
	"sort"
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

func GetWeeksArray(weekStr string) []string {
	if weekStr == "" {
		return GetLastNWeeks(time.Now(), 3*4)
	}
	return strings.Split(weekStr, ",")
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

type ISOWeek struct {
	Year int
	Week int
}

func dowDec31OfYear(y int) int {
	return (y + (y / 4) - (y / 100) + (y / 400)) % 7
}

func isoWeeksInYear(year int) int {
	// If current year ends on Thursday or Previous ends on Wednesday
	if dowDec31OfYear(year) == 4 || dowDec31OfYear(year-1) == 3 {
		return 53
	}

	return 52
}

func ParseISOWeek(iso string) (ISOWeek, error) {
	parts := strings.Split(iso, "-W")
	if len(parts) != 2 {
		return ISOWeek{}, fmt.Errorf("")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return ISOWeek{}, err
	}

	week, err := strconv.Atoi(parts[1])
	if err != nil {
		return ISOWeek{}, err
	}

	if week < 1 || week > isoWeeksInYear(year) {
		return ISOWeek{}, fmt.Errorf("")
	}

	return ISOWeek{Year: year, Week: week}, nil
}

func SortISOWeeks(isoWeeks []string) []string {
	var sortedWeeks []string
	var parsedWeeks []ISOWeek
	mapping := make(map[ISOWeek]string)

	for _, iso := range isoWeeks {
		parsed, err := ParseISOWeek(iso)
		if err != nil {
			continue
		}
		parsedWeeks = append(parsedWeeks, parsed)
		mapping[parsed] = iso
	}

	sort.Slice(parsedWeeks, func(i, j int) bool {
		if parsedWeeks[i].Year == parsedWeeks[j].Year {
			return parsedWeeks[i].Week < parsedWeeks[j].Week
		}
		return parsedWeeks[i].Year < parsedWeeks[j].Year
	})

	for _, week := range parsedWeeks {
		sortedWeeks = append(sortedWeeks, mapping[week])
	}

	return sortedWeeks
}

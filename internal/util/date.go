package util

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

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

func GetWeeksArray(weekStr string) []string {
	if weekStr == "" {
		return GetLastNWeeks(time.Now(), 3*4)
	}
	return strings.Split(weekStr, ",")
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

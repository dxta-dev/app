package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetStartOfWeek(date time.Time) time.Time {
	offset := int(time.Monday - date.Weekday())

	if offset > 0 {
		offset = -6
	}

	startOfWeek := date.AddDate(0, 0, offset)

	startOfWeek = startOfWeek.Truncate(24 * time.Hour)

	return startOfWeek
}

func GetFormattedWeek(date time.Time) string {
	year, week := date.ISOWeek()

	formattedWeek := fmt.Sprintf("%d-W%02d", year, week)

	return formattedWeek
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

	firstDayOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	daysToStartOfWeek := week*7-6
	startOfWeek := firstDayOfYear.AddDate(0, 0, daysToStartOfWeek)

	return startOfWeek, nil
}

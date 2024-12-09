package data

import (
	"strings"
)

type OverallWeeklyData[T any] struct {
	Overall T               `json:"overall"`
	Weekly  []WeeklyData[T] `json:"weekly"`
}

type WeeklyData[T any] struct {
	Week string `json:"week"`
	Data T      `json:"data"`
}

type Statistics struct {
	Average      *float64 `json:"average"`
	Median       *float64 `json:"median"`
	Percentile75 *float64 `json:"percentile75"`
	Percentile95 *float64 `json:"percentile95"`
	Total        *float64 `json:"total"`
	Count        *float64 `json:"count"`
}

type Value struct {
	Value *int `json:"value"`
}

func getTeamSubquery() string {
	return "AND author.external_id in (SELECT member FROM tenant_team_members WHERE team = ?)"
}

func getWeeksPlaceholder(weeksLen int) string {
	return strings.Repeat("?,", weeksLen-1) + "?"
}

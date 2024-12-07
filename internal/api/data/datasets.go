package data

import (
	"database/sql"
	"fmt"
	"strings"
)

type AggregatedStatistics = OverallWeeklyData[Statistics]

type AggregatedValues = OverallWeeklyData[Value]

type WeeklyStatisticsData = WeeklyData[Statistics]

type WeeklyValueData = WeeklyData[Value]

type AggregatedCycleTimeStatistics = OverallWeeklyData[CycleTimeStatistics]

type WeeklyCycleTimeStatistics = WeeklyData[CycleTimeStatistics]

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

type CycleTimeStatistics struct {
	CodingTime Statistics `json:"coding_time"`
	PickupTime Statistics `json:"pickup_time"`
	ReviewTime Statistics `json:"review_time"`
	DeployTime Statistics `json:"deploy_time"`
}

func ScanAggregatedStatsRows(rows *sql.Rows, weeks []string) (*AggregatedStatistics, error) {
	datasetByWeek := make(map[string]WeeklyStatisticsData)
	var aggregated Statistics

	nullWeeksCount := 0
	for rows.Next() {
		var nullableWeek sql.NullString
		var data Statistics
		if err := rows.Scan(
			&nullableWeek,
			&data.Average,
			&data.Median,
			&data.Percentile75,
			&data.Percentile95,
			&data.Total,
			&data.Count,
		); err != nil {
			return nil, err
		}

		if nullableWeek.Valid {
			datasetByWeek[nullableWeek.String] = WeeklyStatisticsData{
				Week: nullableWeek.String,
				Data: data,
			}
		} else {
			aggregated = data
			nullWeeksCount++
		}

		if nullWeeksCount > 1 {
			return nil, fmt.Errorf("ScanAggregatedStatsRows found more than one aggregate row")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var weeklies []WeeklyStatisticsData
	for _, week := range weeks {
		dataPoint, ok := datasetByWeek[week]
		if !ok {
			dataPoint = WeeklyStatisticsData{
				Week: week,
				Data: Statistics{
					Average:      nil,
					Median:       nil,
					Percentile75: nil,
					Percentile95: nil,
					Total:        nil,
					Count:        nil,
				},
			}
		}
		weeklies = append(weeklies, dataPoint)
	}

	return &AggregatedStatistics{
		Overall: aggregated,
		Weekly:  weeklies,
	}, nil
}

func buildQueryAggregatedValues(baseQuery string) string {
	return fmt.Sprintf(`
		WITH dataset AS (%s)
		SELECT NULL as week, SUM(value) AS value FROM dataset
		UNION ALL
		SELECT week, value FROM dataset;`,
		baseQuery,
	)
}

func ScanAggregatedValuesRows(rows *sql.Rows, weeks []string) (*AggregatedValues, error) {
	datasetByWeek := make(map[string]WeeklyValueData)
	var aggregated Value

	nullWeeksCount := 0
	for rows.Next() {
		var nullableWeek sql.NullString
		var data Value
		if err := rows.Scan(
			&nullableWeek,
			&data.Value,
		); err != nil {
			return nil, err
		}

		if nullableWeek.Valid {
			datasetByWeek[nullableWeek.String] = WeeklyValueData{
				Week: nullableWeek.String,
				Data: data,
			}
		} else {
			aggregated = data
			nullWeeksCount++
		}

		if nullWeeksCount > 1 {
			return nil, fmt.Errorf("ScanAggregatedValuesRows found more than one aggregate row")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var weeklies []WeeklyValueData
	for _, week := range weeks {
		dataPoint, ok := datasetByWeek[week]
		if !ok {
			dataPoint = WeeklyValueData{
				Week: week,
				Data: Value{
					Value: nil,
				},
			}
		}
		weeklies = append(weeklies, dataPoint)
	}

	return &AggregatedValues{
		Overall: aggregated,
		Weekly:  weeklies,
	}, nil
}

func getTeamSubquery() string {
	return "AND author.external_id in (SELECT member FROM tenant_team_members WHERE team = ?)"
}

func getWeeksPlaceholder(weeksLen int) string {
	return strings.Repeat("?,", weeksLen-1) + "?"
}

func buildQueryAggregatedStats(baseQuery string) string {
	return fmt.Sprintf(`
		WITH dataset AS (%s),
		data_by_week AS (
			SELECT
				week AS week,
				AVG(value) AS avg,
				MEDIAN(value) AS p50,
				PERCENTILE_75(value) AS p75,
				PERCENTILE_95(value) AS p95,
				SUM(value) as total,
				COUNT(*) as count
			FROM dataset
			GROUP BY week
		),
		data_total AS (
			SELECT AVG(value) as avg,
				MEDIAN(value) as p50,
				PERCENTILE_75(value) as p75,
				PERCENTILE_95(value) as p95,
				SUM(value) as total,
				COUNT(*) as count
			FROM dataset
		)
		SELECT
			NULL as week,
			avg,
			p50,
			p75,
			p95,
			total,
			count
		FROM data_total
		UNION ALL
		SELECT
			week,
			avg,
			p50,
			p75,
			p95,
			total,
			count
		FROM data_by_week;`,
		baseQuery,
	)
}

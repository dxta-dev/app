package api

import (
	"database/sql"
	"fmt"

	"golang.org/x/exp/constraints"
)

func buildQueryAggregatedStatisticData(baseQuery string) string {
	return fmt.Sprintf(`WITH dataset AS (%s),
data_by_week AS (
	SELECT
		x AS WEEK,
		AVG(y) AS AVG,
		MEDIAN(y) AS P50,
		PERCENTILE_75(y) AS P75,
		PERCENTILE_95(y) AS P95,
		SUM(y) as TOTAL,
  	COUNT(y) as COUNT
	FROM dataset
	GROUP BY WEEK
),
data_total AS (
  SELECT AVG(y) as AVG,
  MEDIAN(y) as P50,
  PERCENTILE_75(y) as P75,
  PERCENTILE_95(y) as P95,
  SUM(y) as TOTAL,
  COUNT(y) as COUNT
	FROM dataset
)
SELECT NULL as WEEK, AVG, P50, P75, P95, TOTAL, COUNT FROM data_total
UNION ALL
SELECT WEEK, AVG, P50, P75, P95, TOTAL, COUNT FROM data_by_week;`,
		baseQuery,
	)
}

type AggregatedStats[T constraints.Ordered] struct {
	Overall StatsData[T]         `json:"overall"`
	Weekly  []WeeklyStatsData[T] `json:"weekly"`
}

type WeeklyStatsData[T constraints.Ordered] struct {
	Week string `json:"week"`
	StatsData[T]
}

type StatsData[T constraints.Ordered] struct {
	Average      *T `json:"average"`
	Median       *T `json:"median"`
	Percentile75 *T `json:"percentile75"`
	Percentile95 *T `json:"percentile95"`
	Total        *T `json:"total"`
	Count        *T `json:"count"`
}

func ScanAggregatedStatsRows[T constraints.Ordered](rows *sql.Rows, weeks []string) (*AggregatedStats[T], error) {
	datasetByWeek := make(map[string]WeeklyStatsData[T])
	var aggregated StatsData[T]

	nullWeeksCount := 0
	for rows.Next() {
		var nullableWeek sql.NullString
		var data StatsData[T]
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
			datasetByWeek[nullableWeek.String] = WeeklyStatsData[T]{
				Week:      nullableWeek.String,
				StatsData: data,
			}
		} else {
			aggregated = data
			nullWeeksCount++
		}

		if nullWeeksCount > 1 {
			return nil, fmt.Errorf("ScanAggregatedStatsRows found more than one aggregate rows")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var weeklies []WeeklyStatsData[T]
	for _, week := range weeks {
		dataPoint, ok := datasetByWeek[week]
		if !ok {
			dataPoint = WeeklyStatsData[T]{
				Week: week,
				StatsData: StatsData[T]{
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

	return &AggregatedStats[T]{
		Overall: aggregated,
		Weekly:  weeklies,
	}, nil
}

func buildQueryAggregatedValueData(baseQuery string) string {
	return fmt.Sprintf(`WITH dataset AS (%s)
SELECT NULL as WEEK, SUM(VALUE) AS VALUE FROM dataset
UNION ALL
SELECT WEEK, VALUE FROM dataset;`,
		baseQuery,
	)
}

type AggregatedValues[T constraints.Ordered] struct {
	Overall ValueData[T]         `json:"overall"`
	Weekly  []WeeklyValueData[T] `json:"weekly"`
}

type WeeklyValueData[T constraints.Ordered] struct {
	Week string `json:"week"`
	ValueData[T]
}

type ValueData[T constraints.Ordered] struct {
	Value *T `json:"value"`
}

func ScanAggregatedValuesRows[T constraints.Ordered](rows *sql.Rows, weeks []string) (*AggregatedValues[T], error) {
	datasetByWeek := make(map[string]WeeklyValueData[T])
	var aggregated ValueData[T]

	nullWeeksCount := 0
	for rows.Next() {
		var nullableWeek sql.NullString
		var data ValueData[T]
		if err := rows.Scan(
			&nullableWeek,
			&data.Value,
		); err != nil {
			return nil, err
		}

		if nullableWeek.Valid {
			datasetByWeek[nullableWeek.String] = WeeklyValueData[T]{
				Week:      nullableWeek.String,
				ValueData: data,
			}
		} else {
			aggregated = data
			nullWeeksCount++
		}

		if nullWeeksCount > 1 {
			return nil, fmt.Errorf("ScanAggregatedValueDataRows found more than one aggregate rows")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var weeklies []WeeklyValueData[T]
	for _, week := range weeks {
		dataPoint, ok := datasetByWeek[week]
		if !ok {
			dataPoint = WeeklyValueData[T]{
				Week: week,
				ValueData: ValueData[T]{
					Value: nil,
				},
			}
		}
		weeklies = append(weeklies, dataPoint)
	}

	return &AggregatedValues[T]{
		Overall: aggregated,
		Weekly:  weeklies,
	}, nil
}

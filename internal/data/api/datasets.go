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
	ORDER BY WEEK ASC
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

type AggregatedStatisticData[T constraints.Ordered] struct {
	Aggregated StatisticDataPoint[T]         `json:"aggregated"`
	Weeks      []WeeklyStatisticDataPoint[T] `json:"weeks"`
}

type WeeklyStatisticDataPoint[T constraints.Ordered] struct {
	StatisticDataPoint[T]
	Week string `json:"week"`
}

type StatisticDataPoint[T constraints.Ordered] struct {
	Average      *T `json:"average"`
	Median       *T `json:"median"`
	Percentile75 *T `json:"percentile75"`
	Percentile95 *T `json:"percentile95"`
	Total        *T `json:"total"`
	Count        *T `json:"count"`
}

func ScanAggregatedStatisticDataRows[T constraints.Ordered](rows *sql.Rows, weeks []string) (*AggregatedStatisticData[T], error) {
	datasetByWeek := make(map[string]WeeklyStatisticDataPoint[T])
	var aggregated StatisticDataPoint[T]

	nullWeeksCount := 0
	for rows.Next() {
		var nullableWeek sql.NullString
		var data StatisticDataPoint[T]
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
			datasetByWeek[nullableWeek.String] = WeeklyStatisticDataPoint[T]{
				Week:               nullableWeek.String,
				StatisticDataPoint: data,
			}
		} else {
			aggregated = data
			nullWeeksCount++
		}

		if nullWeeksCount > 1 {
			return nil, fmt.Errorf("ScanAggregatedStatisticDataRows found more than one aggregate rows")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var weeklies []WeeklyStatisticDataPoint[T]
	for _, week := range weeks {
		dataPoint, ok := datasetByWeek[week]
		if !ok {
			dataPoint = WeeklyStatisticDataPoint[T]{
				Week: week,
				StatisticDataPoint: StatisticDataPoint[T]{
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

	return &AggregatedStatisticData[T]{
		Aggregated: aggregated,
		Weeks:      weeklies,
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

type AggregatedValueData[T constraints.Ordered] struct {
	Aggregated ValueDataPoint[T]         `json:"aggregated"`
	Weeks      []WeeklyValueDataPoint[T] `json:"weeks"`
}

type WeeklyValueDataPoint[T constraints.Ordered] struct {
	ValueDataPoint[T]
	Week string `json:"week"`
}

type ValueDataPoint[T constraints.Ordered] struct {
	Value *T `json:"value"`
}

func ScanAggregatedValueDataRows[T constraints.Ordered](rows *sql.Rows, weeks []string) (*AggregatedValueData[T], error) {
	datasetByWeek := make(map[string]WeeklyValueDataPoint[T])
	var aggregated ValueDataPoint[T]

	nullWeeksCount := 0
	for rows.Next() {
		var nullableWeek sql.NullString
		var data ValueDataPoint[T]
		if err := rows.Scan(
			&nullableWeek,
			&data.Value,
		); err != nil {
			return nil, err
		}

		if nullableWeek.Valid {
			datasetByWeek[nullableWeek.String] = WeeklyValueDataPoint[T]{
				Week:           nullableWeek.String,
				ValueDataPoint: data,
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

	var weeklies []WeeklyValueDataPoint[T]
	for _, week := range weeks {
		dataPoint, ok := datasetByWeek[week]
		if !ok {
			dataPoint = WeeklyValueDataPoint[T]{
				Week: week,
				ValueDataPoint: ValueDataPoint[T]{
					Value: nil,
				},
			}
		}
		weeklies = append(weeklies, dataPoint)
	}

	return &AggregatedValueData[T]{
		Aggregated: aggregated,
		Weeks:      weeklies,
	}, nil
}

type ValueData struct {
	Week  string `json:"week"`
	Value *int   `json:"value"`
}

func ScanValueDatasetRows(rows *sql.Rows, weeks []string) ([]ValueData, error) {
	datasetByWeek := make(map[string]ValueData)

	for rows.Next() {
		var dataPoint ValueData
		if err := rows.Scan(
			&dataPoint.Week,
			&dataPoint.Value,
		); err != nil {
			return nil, err
		}

		datasetByWeek[dataPoint.Week] = dataPoint
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var dataset []ValueData
	for _, week := range weeks {
		dataPoint, ok := datasetByWeek[week]
		if !ok {
			dataPoint = ValueData{
				Week:  week,
				Value: nil,
			}
		}
		dataset = append(dataset, dataPoint)
	}

	return dataset, nil
}

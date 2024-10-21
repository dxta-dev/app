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
	var IGNORED sql.Null[any]

	var aggregated StatisticDataPoint[T]
	if rows.Next() {
		if err := rows.Scan(
			&IGNORED,
			&aggregated.Average,
			&aggregated.Median,
			&aggregated.Percentile75,
			&aggregated.Percentile95,
			&aggregated.Total,
			&aggregated.Count,
		); err != nil {
			return nil, err
		}
	}

	for rows.Next() {
		var weekly WeeklyStatisticDataPoint[T]
		if err := rows.Scan(
			&weekly.Week,
			&weekly.Average,
			&weekly.Median,
			&weekly.Percentile75,
			&weekly.Percentile95,
			&weekly.Total,
			&weekly.Count,
		); err != nil {
			return nil, err
		}

		datasetByWeek[weekly.Week] = weekly
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

type StatisticData[T constraints.Ordered] struct {
	Week         string `json:"week"`
	Average      *T     `json:"average"`
	Median       *T     `json:"median"`
	Percentile75 *T     `json:"percentile75"`
	Percentile95 *T     `json:"percentile95"`
	Total        *T     `json:"total"`
	Count        *T     `json:"count"`
}

func ScanStatisticDatasetRows[T constraints.Ordered](rows *sql.Rows, weeks []string) ([]StatisticData[T], error) {
	datasetByWeek := make(map[string]StatisticData[T])

	for rows.Next() {
		var dataPoint StatisticData[T]
		if err := rows.Scan(
			&dataPoint.Week,
			&dataPoint.Average,
			&dataPoint.Median,
			&dataPoint.Percentile75,
			&dataPoint.Percentile95,
			&dataPoint.Total,
			&dataPoint.Count,
		); err != nil {
			return nil, err
		}

		datasetByWeek[dataPoint.Week] = dataPoint
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var dataset []StatisticData[T]
	for _, week := range weeks {
		dataPoint, ok := datasetByWeek[week]
		if !ok {
			dataPoint = StatisticData[T]{
				Week:         week,
				Average:      nil,
				Median:       nil,
				Percentile75: nil,
				Percentile95: nil,
				Total:        nil,
				Count:        nil,
			}
		}
		dataset = append(dataset, dataPoint)
	}

	return dataset, nil
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

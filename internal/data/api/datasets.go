package api

import (
	"database/sql"

	"golang.org/x/exp/constraints"
)

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

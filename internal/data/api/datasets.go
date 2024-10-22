package api

import (
	"database/sql"
	"fmt"
)

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

type AggregatedStats struct {
	Overall StatsData         `json:"overall"`
	Weekly  []WeeklyStatsData `json:"weekly"`
}

type WeeklyStatsData struct {
	Week string `json:"week"`
	StatsData
}

type StatsData struct {
	Average      *float64 `json:"average"`
	Median       *float64 `json:"median"`
	Percentile75 *float64 `json:"percentile75"`
	Percentile95 *float64 `json:"percentile95"`
	Total        *float64 `json:"total"`
	Count        *float64 `json:"count"`
}

func ScanAggregatedStatsRows(rows *sql.Rows, weeks []string) (*AggregatedStats, error) {
	datasetByWeek := make(map[string]WeeklyStatsData)
	var aggregated StatsData

	nullWeeksCount := 0
	for rows.Next() {
		var nullableWeek sql.NullString
		var data StatsData
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
			datasetByWeek[nullableWeek.String] = WeeklyStatsData{
				Week:      nullableWeek.String,
				StatsData: data,
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

	var weeklies []WeeklyStatsData
	for _, week := range weeks {
		dataPoint, ok := datasetByWeek[week]
		if !ok {
			dataPoint = WeeklyStatsData{
				Week: week,
				StatsData: StatsData{
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

	return &AggregatedStats{
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

type AggregatedValues struct {
	Overall ValueData         `json:"overall"`
	Weekly  []WeeklyValueData `json:"weekly"`
}

type WeeklyValueData struct {
	Week string `json:"week"`
	ValueData
}

type ValueData struct {
	Value *int `json:"value"`
}

func ScanAggregatedValuesRows(rows *sql.Rows, weeks []string) (*AggregatedValues, error) {
	datasetByWeek := make(map[string]WeeklyValueData)
	var aggregated ValueData

	nullWeeksCount := 0
	for rows.Next() {
		var nullableWeek sql.NullString
		var data ValueData
		if err := rows.Scan(
			&nullableWeek,
			&data.Value,
		); err != nil {
			return nil, err
		}

		if nullableWeek.Valid {
			datasetByWeek[nullableWeek.String] = WeeklyValueData{
				Week:      nullableWeek.String,
				ValueData: data,
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
				ValueData: ValueData{
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

package data

import (
	"context"
	"database/sql"
	"fmt"
)

type AggregatedStatistics = OverallWeeklyData[Statistics]
type WeeklyStatisticsData = WeeklyData[Statistics]
type AggregatedStatisticsQuery = string

func (d DB) GetAggregatedStatistics(
	ctx context.Context,
	query AggregatedStatisticsQuery,
	namespace string,
	repository string,
	weeks []string,
	team *int64,
) (*AggregatedStatistics, error) {
	queryParamLength := len(weeks)

	queryParams := make([]any, queryParamLength)
	for i, v := range weeks {
		queryParams[i] = v
	}

	queryParams = append(queryParams, namespace)
	queryParams = append(queryParams, repository)

	if team != nil {
		queryParams = append(queryParams, team)
	}

	rows, err := d.db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result, err := ScanAggregatedStatisticsRows(rows, weeks)

	if err != nil {
		return nil, err
	}

	return result, nil
}
func ScanAggregatedStatisticsRows(rows *sql.Rows, weeks []string) (*AggregatedStatistics, error) {
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

func buildQueryAggregatedStatistics(baseQuery string) AggregatedStatisticsQuery {
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

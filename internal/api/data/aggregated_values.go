package data

import (
	"context"
	"database/sql"
	"fmt"
)

type AggregatedValues = OverallWeeklyData[Value]
type WeeklyValueData = WeeklyData[Value]
type AggregatedValuesQuery = string

func (d DB) GetAggregatedValues(ctx context.Context, query AggregatedValuesQuery, namespace string, repository string, weeks []string, team *int64) (*AggregatedValues, error) {
	queryParamLength := len(weeks)

	queryParams := make([]interface{}, queryParamLength)
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
	result, err := scanAggregatedValuesRows(rows, weeks)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func scanAggregatedValuesRows(rows *sql.Rows, weeks []string) (*AggregatedValues, error) {
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

func buildQueryAggregatedValues(baseQuery string) AggregatedValuesQuery {
	return fmt.Sprintf(`
		WITH dataset AS (%s)
		SELECT NULL as week, SUM(value) AS value FROM dataset
		UNION ALL
		SELECT week, value FROM dataset;`,
		baseQuery,
	)
}

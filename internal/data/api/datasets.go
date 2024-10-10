package api

import "database/sql"

type ValueIntegerDataset struct {
	Week  string `json:"week"`
	Value *int   `json:"value"`
}
type ValueRealDataset struct {
	Week  string  `json:"week"`
	Value float64 `json:"value"`
}

type StatisticIntegerDataset struct {
	Week         string `json:"week"`
	Average      int    `json:"average"`
	Median       int    `json:"median"`
	Percentile75 int    `json:"percentile75"`
	Percentile95 int    `json:"percentile95"`
}
type StatisticRealDataset struct {
	Week         string  `json:"week"`
	Average      float64 `json:"average"`
	Median       float64 `json:"median"`
	Percentile75 float64 `json:"percentile75"`
	Percentile95 float64 `json:"percentile95"`
}

type CountIntegerDataset struct {
	Week  string `json:"week"`
	Count int    `json:"value"`
}
type CountRealDataset struct {
	Week  string  `json:"week"`
	Count float64 `json:"value"`
}

func ScanValueIntegerDatasetRows(rows *sql.Rows, weeks []string) ([]ValueIntegerDataset, error) {
	datasetByWeek := make(map[string]ValueIntegerDataset)

	for rows.Next() {
		var dataPoint ValueIntegerDataset
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

	var dataset []ValueIntegerDataset
	for _, week := range weeks {
		dataPoint, ok := datasetByWeek[week]
		if !ok {
			dataPoint = ValueIntegerDataset{
				Week:  week,
				Value: nil,
			}
		}
		dataset = append(dataset, dataPoint)
	}

	return dataset, nil
}

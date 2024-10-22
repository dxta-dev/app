package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func GetDeployFrequency(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string) (*AggregatedValues, error) {

	queryParamLength := len(weeks)

	weeksPlaceholder := strings.Repeat("?,", len(weeks)-1) + "?"

	queryParams := make([]interface{}, queryParamLength)
	for i, v := range weeks {
		queryParams[i] = v
	}

	queryParams = append(queryParams, namespace)
	queryParams = append(queryParams, repository)

	query := buildQueryAggregatedValueData(fmt.Sprintf(`
	SELECT
		deploy_dates.week AS week,
		COUNT(*) AS value
	FROM transform_deployments AS deploy
	JOIN transform_dates AS deploy_dates
	ON deploy.deployed_at = deploy_dates.id
	JOIN transform_repositories AS repo
	ON deploy.repository_id = repo.id
	WHERE deploy_dates.week IN (%s)
	AND repo.namespace_name = ?
	AND repo.name = ?
	GROUP BY deploy_dates.week`,
		weeksPlaceholder,
	))

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	deployFrequencies, err := ScanAggregatedValuesRows(rows, weeks)

	if err != nil {
		return nil, err
	}

	return deployFrequencies, nil
}

package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type DeployFrequency struct {
	Week  string  `json:"week"`
	Value float32 `json:"value"`
}

/*

	SELECT
		deploy_dates.week,
		CAST(COUNT (*) AS REAL) / 7
	FROM transform_merge_request_metrics AS metrics
	JOIN transform_repositories AS repo
	ON repo.id = metrics.repository
	JOIN transform_merge_request_fact_dates_junk AS dates_junk
	ON metrics.dates_junk = dates_junk.id
	JOIN transform_dates AS deploy_dates
	ON dates_junk.deployed_at = deploy_dates.id
	JOIN transform_merge_request_fact_users_junk AS uj
	ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
	ON uj.author = author.id
	WHERE deploy_dates.week IN ("2024-W26", "2024-W27", "2024-W28", "2024-W29", "2024-W30", "2024-W31", "2024-W32", "2024-W33", "2024-W34", "2024-W35", "2024-W36", "2024-W37")
	AND repo.namespace_name = "calcom"
	AND repo.name = "cal.com"
	AND author.external_id in (SELECT member FROM tenant_team_members WHERE team = 1)
	AND author.bot = 0
	GROUP BY deploy_dates.week
	ORDER BY deploy_dates.week ASC;

*/

func GetDeployFrequency(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) ([]DeployFrequency, error) {

	teamQuery := ""
	queryParamLength := len(weeks)

	if team != nil {
		teamQuery = "AND author.external_id in (SELECT member FROM tenant_team_members WHERE team = ?)"
		queryParamLength += 1
	}

	weeksPlaceholder := strings.Repeat("?,", len(weeks)-1) + "?"

	queryParams := make([]interface{}, queryParamLength)
	for i, v := range weeks {
		queryParams[i] = v
	}

	queryParams = append(queryParams, namespace)
	queryParams = append(queryParams, repository)

	if team != nil {
		queryParams = append(queryParams, team)
	}

	query := fmt.Sprintf(`
	SELECT
		deploy_dates.week,
		CAST(COUNT (*) AS REAL) / 7
	FROM transform_merge_request_metrics AS metrics
	JOIN transform_repositories AS repo
	ON repo.id = metrics.repository
	JOIN transform_merge_request_fact_dates_junk AS dates_junk
	ON metrics.dates_junk = dates_junk.id
	JOIN transform_dates AS deploy_dates
	ON dates_junk.deployed_at = deploy_dates.id
	JOIN transform_merge_request_fact_users_junk AS uj
	ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
	ON uj.author = author.id
	WHERE deploy_dates.week IN (%s)
	AND repo.namespace_name = ?
	AND repo.name = ?
	%s
	AND author.bot = 0
	GROUP BY deploy_dates.week
	ORDER BY deploy_dates.week ASC;
	`,
		weeksPlaceholder,
		teamQuery,
	)

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var deployFrequencies []DeployFrequency

	for rows.Next() {
		var deployFrequency DeployFrequency

		if err := rows.Scan(
			&deployFrequency.Week,
			&deployFrequency.Value,
		); err != nil {
			return nil, err
		}
		deployFrequencies = append(deployFrequencies, deployFrequency)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return deployFrequencies, nil
}

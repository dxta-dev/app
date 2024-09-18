package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type CycleTime struct {
	Week         string `json:"week"`
	Average      int    `json:"average"`
	Median       int    `json:"median"`
	Percentile75 int    `json:"percentile75"`
	Percentile95 int    `json:"percentile95"`
}

/*
	SELECT
		deployedAt.week as WEEK,
		FLOOR(AVG(metrics.coding_duration + metrics.review_start_delay + metrics.review_duration + metrics.deploy_duration)) AS AVG,
		FLOOR(MEDIAN(metrics.coding_duration + metrics.review_start_delay + metrics.review_duration + metrics.deploy_duration)) as P50,
		FLOOR(PERCENTILE_75(metrics.coding_duration + metrics.review_start_delay + metrics.review_duration + metrics.deploy_duration)) as P75,
		FLOOR(PERCENTILE_95(metrics.coding_duration + metrics.review_start_delay + metrics.review_duration + metrics.deploy_duration)) as P95
		FROM transform_merge_request_metrics AS metrics
	JOIN transform_repositories AS repo
		ON repo.id = metrics.repository
	JOIN transform_merge_request_fact_dates_junk AS dj
		ON metrics.dates_junk = dj.id
	JOIN transform_dates AS deployedAt
		ON dj.deployed_at = deployedAt.id
	JOIN transform_merge_request_fact_users_junk AS uj
		ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
		ON uj.author = author.id
	WHERE deployedAt.week IN ("2024-W26", "2024-W27", "2024-W28", "2024-W29", "2024-W30", "2024-W31", "2024-W32", "2024-W33", "2024-W34", "2024-W35", "2024-W36", "2024-W37")
	AND repo.name = "cal.com"
	AND repo.namespace_name = "calcom"
	AND author.external_id IN (SELECT member FROM tenant_team_members WHERE team = 1)
	AND author.bot = 0
	GROUP BY deployedAt.week
	ORDER BY deployedAt.week ASC;
*/

func GetCycleTime(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) ([]CycleTime, error) {

	teamQuery := ""
	queryParamLength := len(weeks)

	if team != nil {
		teamQuery = "AND author.external_id in (SELECT member FROM tenant_team_members WHERE team = ?)"
		queryParamLength += 1
	}

	weeksPlaceholder := strings.Repeat("?,", queryParamLength)

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
			deployedAt.week as WEEK,
			FLOOR(AVG(metrics.coding_duration + metrics.review_start_delay + metrics.review_duration + metrics.deploy_duration)) AS AVG,
			FLOOR(MEDIAN(metrics.coding_duration + metrics.review_start_delay + metrics.review_duration + metrics.deploy_duration)) as P50,
			FLOOR(PERCENTILE_75(metrics.coding_duration + metrics.review_start_delay + metrics.review_duration + metrics.deploy_duration)) as P75,
			FLOOR(PERCENTILE_95(metrics.coding_duration + metrics.review_start_delay + metrics.review_duration + metrics.deploy_duration)) as P95
			FROM transform_merge_request_metrics AS metrics
		JOIN transform_repositories AS repo
			ON repo.id = metrics.repository
		JOIN transform_merge_request_fact_dates_junk AS dj
			ON metrics.dates_junk = dj.id
		JOIN transform_dates AS deployedAt
			ON dj.deployed_at = deployedAt.id
		JOIN transform_merge_request_fact_users_junk AS uj
			ON metrics.users_junk = uj.id
		JOIN transform_forge_users AS author
		ON uj.author = author.id
		WHERE deployedAt.week IN (%s)
		AND repo.namespace_name = ?
		AND repo.name = ?
		%s
		AND author.bot = 0
		GROUP BY deployedAt.week
		ORDER BY deployedAt.week ASC;
	`,
		weeksPlaceholder,
		teamQuery,
	)

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var cycleTimes []CycleTime

	for rows.Next() {
		var cycleTime CycleTime

		if err := rows.Scan(
			&cycleTime.Week,
			&cycleTime.Average,
			&cycleTime.Median,
			&cycleTime.Percentile75,
			&cycleTime.Percentile95,
		); err != nil {
			return nil, err
		}
		cycleTimes = append(cycleTimes, cycleTime)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cycleTimes, nil
}

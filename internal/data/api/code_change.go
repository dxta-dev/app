package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type CodeChange = *AggregatedValueData[int]

/*
	SELECT
		dates.week AS WEEK,
		SUM(metrics.mr_size) AS VALUE
	FROM transform_merge_request_metrics AS metrics
	JOIN transform_repositories AS repo
	ON repo.id = metrics.repository
	JOIN transform_merge_request_fact_dates_junk AS dates_junk
	ON metrics.dates_junk = dates_junk.id
	JOIN transform_dates AS dates
	ON dates_junk.merged_at = dates.id
	JOIN transform_merge_request_fact_users_junk AS uj
	ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
	ON uj.author = author.id
	WHERE dates.week IN ("2024-W26", "2024-W27", "2024-W28", "2024-W29", "2024-W30", "2024-W31", "2024-W32", "2024-W33", "2024-W34", "2024-W35", "2024-W36", "2024-W37")
	AND repo.name = "cal.com"
	AND repo.namespace_name = "calcom"
	AND author.external_id in (SELECT member FROM tenant_team_members WHERE team = 1)
	GROUP BY dates.week
	ORDER BY dates.week ASC;
*/

func GetCodeChanges(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) (CodeChange, error) {

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

	query := buildQueryAggregatedValueData(fmt.Sprintf(`
	SELECT
		dates.week AS WEEK,
		SUM(metrics.mr_size) AS VALUE
	FROM transform_merge_request_metrics AS metrics
	JOIN transform_repositories AS repo
		ON repo.id = metrics.repository
	JOIN transform_merge_request_fact_dates_junk AS dates_junk
		ON metrics.dates_junk = dates_junk.id
	JOIN transform_dates AS dates
		ON dates_junk.merged_at = dates.id
	JOIN transform_merge_request_fact_users_junk AS uj
		ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
		ON uj.author = author.id
	JOIN transform_merge_requests AS mrs
		ON metrics.merge_request = mrs.id
	JOIN transform_branches AS branch
		ON mrs.target_branch = branch.id	
	WHERE dates.week IN (%s)
	AND repo.namespace_name = ?
	AND repo.name = ?
	AND branch.name = 'main'
		%s
	AND author.bot = 0
	GROUP BY dates.week`,
		weeksPlaceholder,
		teamQuery,
	))

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	codeChanges, err := ScanAggregatedValueDataRows[int](rows, weeks)

	if err != nil {
		return nil, err
	}

	return codeChanges, nil
}

package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func GetMRPickupTime(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) (*AggregatedStats, error) {

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

	query := buildQueryAggregatedStats(fmt.Sprintf(`
	SELECT
		mergedAt.week AS week,
		metrics.review_start_delay AS value
	FROM transform_merge_request_metrics AS metrics
	JOIN transform_repositories AS repo
		ON repo.id = metrics.repository
	JOIN transform_merge_request_fact_dates_junk AS dj
		ON metrics.dates_junk = dj.id
	JOIN transform_dates AS mergedAt
		ON dj.merged_at = mergedAt.id
	JOIN transform_merge_request_fact_users_junk AS uj
		ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
		ON uj.author = author.id
	JOIN transform_merge_requests AS mrs
		ON metrics.merge_request = mrs.id
	JOIN transform_branches AS branch
		ON mrs.target_branch = branch.id
	WHERE mergedAt.week IN (%s)
	AND repo.namespace_name = ?
	AND repo.name = ?
	AND branch.name = 'main'
	%s
	AND author.bot = 0`,
		weeksPlaceholder,
		teamQuery,
	))

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mrsPickupTime, err := ScanAggregatedStatsRows(rows, weeks)

	if err != nil {
		return nil, err
	}

	return mrsPickupTime, nil
}

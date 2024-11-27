package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func GetDeployTime(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) (*AggregatedStats, error) {

	teamQuery := ""
	queryParamLength := len(weeks)

	if team != nil {
		teamQuery = "AND author.external_id in (SELECT member FROM tenant_team_members WHERE team = ?)"
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
	year, month, day := time.Now().Date()

	currentDate := fmt.Sprintf("%04d-%02d-%02d", year, int(month), day)

	query := buildQueryAggregatedStats(fmt.Sprintf(`
	SELECT
		merged_at.week AS week,
		CASE
		WHEN metrics.deploy_duration = 0
		 THEN (julianday('%s') - julianday(
      CONCAT(dates.year, '-', LPAD(dates.month, 2, '0'), '-', LPAD(dates.day, 2, '0'))
    )) * 86400000
		ELSE metrics.deploy_duration END AS value
		FROM transform_merge_request_metrics AS metrics
	JOIN transform_repositories AS repo
		ON repo.id = metrics.repository
	JOIN transform_merge_request_fact_dates_junk AS dj
		ON metrics.dates_junk = dj.id
	JOIN transform_dates AS merged_at
		ON dj.merged_at = merged_at.id
		JOIN transform_dates AS dates
        ON dj.merged_at = dates.id
	JOIN transform_merge_request_fact_users_junk AS uj
		ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
		ON uj.author = author.id
	JOIN transform_merge_requests AS mrs
		ON metrics.merge_request = mrs.id
	JOIN transform_branches AS branch
		ON mrs.target_branch = branch.id
	WHERE merged_at.week IN (%s)
	AND repo.namespace_name = ?
	AND repo.name = ?
	AND branch.id = repo.default_branch
	%s
	AND author.bot = 0`,
		currentDate,
		weeksPlaceholder,
		teamQuery,
	))

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	deployTimes, err := ScanAggregatedStatsRows(rows, weeks)

	if err != nil {
		return nil, err
	}

	return deployTimes, nil
}

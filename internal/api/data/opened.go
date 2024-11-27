package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func GetMRsOpened(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) (*AggregatedValues, error) {

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

	query := buildQueryAggregatedValues(fmt.Sprintf(`
	SELECT
		opened_dates.week AS week,
		COUNT(*) AS value
	FROM transform_merge_request_metrics AS metrics
	JOIN transform_repositories AS repo
	ON repo.id = metrics.repository
	JOIN transform_merge_request_fact_dates_junk AS dates_junk
	ON metrics.dates_junk = dates_junk.id
	JOIN transform_dates AS opened_dates
	ON dates_junk.opened_at = opened_dates.id
	JOIN transform_merge_request_fact_users_junk AS uj
	ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
	ON uj.author = author.id
	JOIN transform_merge_requests AS mrs
	ON metrics.merge_request = mrs.id
	JOIN transform_branches AS branch
	ON mrs.target_branch = branch.id
	WHERE opened_dates.week IN (%s)
	AND repo.namespace_name = ?
	AND repo.name = ?
	AND branch.id = repo.default_branch
	%s
	AND author.bot = 0
	GROUP BY opened_dates.week`,
		weeksPlaceholder,
		teamQuery,
	))

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mrsOpened, err := ScanAggregatedValuesRows(rows, weeks)

	if err != nil {
		return nil, err
	}

	return mrsOpened, nil
}

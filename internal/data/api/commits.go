package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func GetCommits(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) (*AggregatedValues, error) {

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
		commitedAt.week AS WEEK,
		COUNT(*) AS VALUE
	FROM transform_merge_request_events AS ev
	JOIN transform_dates AS commitedAt
	ON commitedAt.id = ev.commited_at
	JOIN transform_repositories AS repo
	ON repo.id = ev.repository
	JOIN transform_forge_users AS author
	ON ev.actor = author.id
	JOIN transform_merge_requests AS mrs
	ON ev.merge_request = mrs.id
	JOIN transform_branches AS branch
	ON mrs.target_branch = branch.id
	WHERE commitedAt.week IN (%s)
	AND ev.merge_request_event_type = 9
	AND repo.namespace_name = ?
	AND repo.name = ?
	AND branch.name = 'main'
	%s
	AND author.bot = 0
	GROUP BY commitedAt.week`,
		weeksPlaceholder,
		teamQuery,
	))

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	commits, err := ScanAggregatedValuesRows(rows, weeks)

	if err != nil {
		return nil, err
	}

	return commits, nil
}

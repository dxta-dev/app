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
		committedAt.week AS week,
		COUNT(*) AS value
	FROM transform_merge_request_events AS ev
	JOIN transform_dates AS committedAt
	ON committedAt.id = ev.commited_at
	JOIN transform_repositories AS repo
	ON repo.id = ev.repository
	JOIN transform_forge_users AS author
	ON ev.actor = author.id
	JOIN transform_merge_requests AS mrs
	ON ev.merge_request = mrs.id
	JOIN transform_branches AS branch
	ON mrs.target_branch = branch.id
	WHERE committedAt.week IN (%s)
	AND ev.merge_request_event_type = 9
	AND repo.namespace_name = ?
	AND repo.name = ?
	AND branch.name = 'main'
	%s
	AND author.bot = 0
	GROUP BY committedAt.week`,
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

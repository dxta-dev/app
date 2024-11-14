package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func GetTotalReviews(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) (*AggregatedValues, error) {
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
		occurredAt.week AS week,
		COUNT(*) AS value
	FROM transform_merge_request_events AS ev
	JOIN transform_repositories AS repo
	ON repo.id = ev.repository
	JOIN transform_dates AS occurredAt
	ON occurredAt.id = ev.occured_on
	JOIN transform_forge_users AS author
	ON ev.actor = author.id
	JOIN transform_merge_requests AS mrs
	ON ev.merge_request = mrs.id
	JOIN transform_branches AS branch
	ON mrs.target_branch = branch.id
	WHERE ev.merge_request_event_type = 15
	AND occurredAt.week IN (%s)
	AND repo.namespace_name = ?
	AND repo.name = ?
	AND branch.id = repo.default_branch
	%s
	AND author.bot = 0
	GROUP BY occurredAt.week`,
		weeksPlaceholder,
		teamQuery,
	))

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	totalReviews, err := ScanAggregatedValuesRows(rows, weeks)

	if err != nil {
		return nil, err
	}

	return totalReviews, nil
}

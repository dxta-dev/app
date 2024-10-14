package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Commits = ValueData

/*
	SELECT
		commitedAt.week,
		COUNT (*)
	FROM transform_merge_request_events AS ev
	JOIN transform_dates AS commitedAt
	ON commitedAt.id = ev.commited_at
	JOIN transform_repositories AS repo
	ON repo.id = ev.repository
	JOIN transform_forge_users AS author
	ON ev.actor = author.id
	WHERE commitedAt.week IN ("2024-W26", "2024-W27", "2024-W28", "2024-W29", "2024-W30", "2024-W31", "2024-W32", "2024-W33", "2024-W34", "2024-W35", "2024-W36", "2024-W37")
	AND ev.merge_request_event_type = 9
	AND repo.namespace_name = "calcom"
	AND repo.name = "cal.com"
	AND author.external_id in (SELECT member FROM tenant_team_members WHERE team = 1)
	AND author.bot = 0
	GROUP BY commitedAt.week
	ORDER BY commitedAt.week ASC;

*/

func GetCommits(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) ([]Commits, error) {

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
		commitedAt.week,
		COUNT (*)
	FROM transform_merge_request_events AS ev
	JOIN transform_dates AS commitedAt
	ON commitedAt.id = ev.commited_at
	JOIN transform_repositories AS repo
	ON repo.id = ev.repository
	JOIN transform_forge_users AS author
	ON ev.actor = author.id
	WHERE commitedAt.week IN (%s)
	AND ev.merge_request_event_type = 9
	AND repo.namespace_name = ?
	AND repo.name = ?
	%s
	AND author.bot = 0
	GROUP BY commitedAt.week
	ORDER BY commitedAt.week ASC;
	`,
		weeksPlaceholder,
		teamQuery,
	)

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	commits, err := ScanValueDatasetRows(rows, weeks)

	if err != nil {
		return nil, err
	}

	return commits, nil
}

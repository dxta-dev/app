package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type TotalReview = *AggregatedValues

/*

	SELECT
		occuredAt.week,
		COUNT (*)
	FROM transform_merge_request_events AS ev
	JOIN transform_repositories AS repo
	ON repo.id = ev.repository
	JOIN transform_dates AS occuredAt
	ON occuredAt.id = ev.occured_on
	JOIN transform_forge_users AS author
	ON ev.actor = author.id
	WHERE ev.merge_request_event_type = 15
	AND occuredAt.week IN ("2024-W26", "2024-W27", "2024-W28", "2024-W29", "2024-W30", "2024-W31", "2024-W32", "2024-W33", "2024-W34", "2024-W35", "2024-W36", "2024-W37")
	AND repo.namespace_name = "calcom"
	AND repo.name = "cal.com"
	AND author.external_id in (SELECT member FROM tenant_team_members WHERE team = 1)
	AND author.bot = 0
	GROUP BY occuredAt.week
	ORDER BY occuredAt.week ASC;



*/

func GetTotalReviews(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) (TotalReview, error) {

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
		occuredAt.week AS WEEK,
		COUNT(*) AS VALUE
	FROM transform_merge_request_events AS ev
	JOIN transform_repositories AS repo
	ON repo.id = ev.repository
	JOIN transform_dates AS occuredAt
	ON occuredAt.id = ev.occured_on
	JOIN transform_forge_users AS author
	ON ev.actor = author.id
	JOIN transform_merge_requests AS mrs
	ON ev.merge_request = mrs.id
	JOIN transform_branches AS branch
	ON mrs.target_branch = branch.id	
	WHERE ev.merge_request_event_type = 15
	AND occuredAt.week IN (%s)
	AND repo.namespace_name = ?
	AND repo.name = ?
	AND branch.name = 'main'
	%s
	AND author.bot = 0
	GROUP BY occuredAt.week`,
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

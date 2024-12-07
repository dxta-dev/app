package data

import (
	"context"
	"database/sql"
	"fmt"
)

type CodeChanges struct {
	DB *sql.DB
}

func (c CodeChanges) BuildQuery(weeks []string, team *int64) string {
	teamQuery := ""

	if team != nil {
		teamQuery = getTeamSubquery()
	}

	weeksPlaceholder := getWeeksPlaceholder(len(weeks))

	query := `
	SELECT
		dates.week AS week,
		SUM(metrics.mr_size) AS value
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
	AND branch.id = repo.default_branch
		%s
	AND author.bot = 0
	GROUP BY dates.week`

	return buildQueryAggregatedValues(fmt.Sprintf(query, weeksPlaceholder, teamQuery))
}

func (c CodeChanges) GetData(ctx context.Context, namespace string, repository string, weeks []string, team *int64) (*AggregatedValues, error) {
	queryParamLength := len(weeks)

	queryParams := make([]interface{}, queryParamLength)
	for i, v := range weeks {
		queryParams[i] = v
	}

	queryParams = append(queryParams, namespace)
	queryParams = append(queryParams, repository)

	if team != nil {
		queryParams = append(queryParams, team)
	}

	query := c.BuildQuery(weeks, team)

	rows, err := c.DB.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	codeChanges, err := ScanAggregatedValuesRows(rows, weeks)

	if err != nil {
		return nil, err
	}

	return codeChanges, nil
}

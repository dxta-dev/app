package data

import (
	"fmt"
)

func BuildSmallMRsQuery(weeks []string, team *int64) AggregatedValuesQuery {
	teamQuery := ""

	if team != nil {
		teamQuery = getTeamSubquery()
	}

	weeksPlaceholder := getWeeksPlaceholder(len(weeks))

	return buildQueryAggregatedValues(fmt.Sprintf(`
	SELECT
		mergedAt.week AS week,
		COUNT(*) AS value
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
	AND metrics.mr_size <= 250
	AND repo.namespace_name = ?
	AND repo.name = ?
	AND branch.id = repo.default_branch
    %s
	AND author.bot = 0
	GROUP BY mergedAt.week`,
		weeksPlaceholder,
		teamQuery,
	))
}

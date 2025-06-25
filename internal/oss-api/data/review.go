package data

import (
	"fmt"
)

func BuildReviewQuery(weeks []string, team *int64) AggregatedValuesQuery {
	teamQuery := ""

	if team != nil {
		teamQuery = getTeamSubquery()
	}

	weeksPlaceholder := getWeeksPlaceholder(len(weeks))

	return buildQueryAggregatedValues(fmt.Sprintf(`
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
}

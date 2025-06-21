package data

import (
	"fmt"
)

func BuildCommitsQuery(weeks []string, team *int64) AggregatedValuesQuery {
	teamQuery := ""

	if team != nil {
		teamQuery = getTeamSubquery()
	}

	weeksPlaceholder := getWeeksPlaceholder(len(weeks))

	return buildQueryAggregatedValues(fmt.Sprintf(`
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
	AND branch.id = repo.default_branch
	%s
	AND author.bot = 0
	GROUP BY committedAt.week`,
		weeksPlaceholder,
		teamQuery,
	))
}

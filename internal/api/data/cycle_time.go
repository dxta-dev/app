package data

import (
	"fmt"
)

func BuildCycleTimeQuery(weeks []string, team *int64) AggregatedStatisticsQuery {
	teamQuery := ""

	if team != nil {
		teamQuery = getTeamSubquery()
	}

	weeksPlaceholder := getWeeksPlaceholder(len(weeks))

	return fmt.Sprintf(`
	WITH dataset AS (
		SELECT
			mergedAt.week AS week,
			metrics.coding_duration + metrics.review_start_delay + metrics.review_duration + metrics.deploy_duration AS value
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
		AND repo.namespace_name = ?
		AND repo.name = ?
		AND branch.id = repo.default_branch
		%s
		AND author.bot = 0
	),
	data_by_week AS (
		SELECT
			week AS week,
			AVG(value) AS avg,
			MEDIAN(value) AS p50,
			PERCENTILE_75(value) AS p75,
			PERCENTILE_95(value) AS p95,
			SUM(value) as total,
			COUNT(*) as count
		FROM dataset
		GROUP BY week
	),
	data_total AS (
		SELECT AVG(value) as avg,
			MEDIAN(value) as p50,
			PERCENTILE_75(value) as p75,
			PERCENTILE_95(value) as p95,
			SUM(value) as total,
			COUNT(*) as count
		FROM dataset
	)
	SELECT
		NULL as week,
		avg,
		p50,
		p75,
		p95,
		total,
		count
	FROM data_total
	UNION ALL
	SELECT
		week,
		avg,
		p50,
		p75,
		p95,
		total,
		count
	FROM data_by_week;`,
		weeksPlaceholder,
		teamQuery,
	)
}

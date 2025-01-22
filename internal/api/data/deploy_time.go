package data

import (
	"fmt"
)

func BuildDeployTimeQuery(weeks []string, team *int64) AggregatedStatisticsQuery {
	teamQuery := ""

	if team != nil {
		teamQuery = getTeamSubquery()
	}

	weeksPlaceholder := getWeeksPlaceholder(len(weeks))

	return fmt.Sprintf(`
	WITH has_deployment AS (
		SELECT DISTINCT repository_external_id, forge_type
		FROM tenant_deployment_environments
		UNION
		SELECT DISTINCT repository_external_id, forge_type
		FROM tenant_cicd_deploy_workflows
	),
	dataset AS (
		SELECT
			merged_at.week AS week,
			CASE
				WHEN metrics.deploy_duration = 0 THEN
						(unixepoch(date('now')) - unixepoch(
							CONCAT(dates.year, '-', LPAD(dates.month, 2, '0'), '-', LPAD(dates.day, 2, '0'))
						)) * 1000
				ELSE metrics.deploy_duration
			END AS value
			FROM transform_merge_request_metrics AS metrics
		JOIN transform_repositories AS repo
			ON repo.id = metrics.repository
		JOIN has_deployment
			ON has_deployment.repository_external_id = repo.external_id AND has_deployment.forge_type = repo.forge_type
		JOIN transform_merge_request_fact_dates_junk AS dj
			ON metrics.dates_junk = dj.id
		JOIN transform_dates AS merged_at
			ON dj.merged_at = merged_at.id
			JOIN transform_dates AS dates
					ON dj.merged_at = dates.id
		JOIN transform_merge_request_fact_users_junk AS uj
			ON metrics.users_junk = uj.id
		JOIN transform_forge_users AS author
			ON uj.author = author.id
		JOIN transform_merge_requests AS mrs
			ON metrics.merge_request = mrs.id
		JOIN transform_branches AS branch
			ON mrs.target_branch = branch.id
		WHERE merged_at.week IN (%s)
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

package data

import (
	"fmt"
)

func BuildDeployFrequencyQuery(weeks []string) AggregatedValuesQuery {

	weeksPlaceholder := getWeeksPlaceholder(len(weeks))

	return buildQueryAggregatedValues(fmt.Sprintf(`
	SELECT
		deploy_dates.week AS week,
		COUNT(*) AS value
	FROM transform_deployments AS deploy
	JOIN transform_dates AS deploy_dates
	ON deploy.deployed_at = deploy_dates.id
	JOIN transform_repositories AS repo
	ON deploy.repository_id = repo.id
	WHERE deploy_dates.week IN (%s)
	AND repo.namespace_name = ?
	AND repo.name = ?
	GROUP BY deploy_dates.week`,
		weeksPlaceholder,
	))
}

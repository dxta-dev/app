package api

import (
	"database/sql"
	"fmt"
	"strings"
)

type MRsMergedWithoutReview struct {
	Week         string
	Average      int
	Median       int
	Percentile75 int
	Percentile95 int
	Count        int
}

/*
	SELECT
		COUNT(*),
		mergedAt.week
	FROM transform_merge_request_metrics AS metrics
	JOIN transform_merge_request_fact_dates_junk AS dj
	ON metrics.dates_junk = dj.id
	JOIN transform_dates AS mergedAt
	ON dj.merged_at = mergedAt.id
	JOIN transform_merge_request_fact_users_junk AS uj
	ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
	ON uj.author = author.id
	WHERE mergedAt.week IN (%s) and metrics.review_depth = 0
	AND author.bot = 0
	%s
	GROUP BY mergedAt.week;
*/

func GetMRsMergedWithoutReview(db *sql.DB, namespace string, repository string, weeks []string, team *int64) (map[string]MRSize, error) {

	teamQuery := ""
	queryParamLength := len(weeks) + 1 /* repository name */ + 1 /* repository namespace */

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

	_ = fmt.Sprintf(`
	SELECT
		mergedAt.week as WEEK,
		FLOOR(AVG(metrics.mr_size)) as AVG,
		FLOOR(MEDIAN(metrics.mr_size)) AS P50,
		FLOOR(PERCENTILE_75(metrics.mr_size)) AS P75,
		FLOOR(PERCENTILE_95(metrics.mr_size)) as P95,
		COUNT(*) AS C
	FROM transform_merge_request_metrics AS metrics
	JOIN transform_merge_request_fact_dates_junk AS dj
	ON metrics.dates_junk = dj.id
	JOIN transform_dates AS mergedAt
	ON dj.merged_at = mergedAt.id
	JOIN transform_merge_request_fact_users_junk AS uj
	ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
	ON uj.author = author.id
	WHERE mergedAt.week IN (%s)
	AND repo.namespace_name = ?
	AND repo.name = ?
	%s
	AND author.bot = 0
	GROUP BY mergedAt.week;
	`,
		weeksPlaceholder,
		teamQuery,
	)

	return nil, nil
}

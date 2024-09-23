package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type MRPickupTime struct {
	Week         string  `json:"week"`
	Average      float64 `json:"average"`
	Median       float64 `json:"median"`
	Percentile75 float64 `json:"percentile75"`
	Percentile95 float64 `json:"percentile95"`
}

/*
	SELECT
		mergedAt.week as WEEK,
		FLOOR(AVG(metrics.review_start_delay)) AS AVG,
		FLOOR(MEDIAN(metrics.review_start_delay)) as P50,
		FLOOR(PERCENTILE_75(metrics.review_start_delay)) as P75,
		FLOOR(PERCENTILE_95(metrics.review_start_delay)) as P95
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
	WHERE mergedAt.week IN ("2024-W26", "2024-W27", "2024-W28", "2024-W29", "2024-W30", "2024-W31", "2024-W32", "2024-W33", "2024-W34", "2024-W35", "2024-W36", "2024-W37")
	AND repo.name = "cal.com"
	AND repo.namespace_name = "calcom"
	AND author.external_id IN (SELECT member FROM tenant_team_members WHERE team = 1)
	AND author.bot = 0
	GROUP BY mergedAt.week
	ORDER BY mergedAt.week ASC;

*/

func GetMRPickupTime(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) ([]MRPickupTime, error) {

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
		mergedAt.week as WEEK,
		FLOOR(AVG(metrics.review_start_delay)) AS AVG,
		FLOOR(MEDIAN(metrics.review_start_delay)) as P50,
		FLOOR(PERCENTILE_75(metrics.review_start_delay)) as P75,
		FLOOR(PERCENTILE_95(metrics.review_start_delay)) as P95
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
	WHERE mergedAt.week IN (%s)
	AND repo.namespace_name = ?
	AND repo.name = ?
	%s
	AND author.bot = 0
	GROUP BY mergedAt.week
	ORDER BY mergedAt.week ASC;
	`,
		weeksPlaceholder,
		teamQuery,
	)

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var mrsPickupTime []MRPickupTime

	for rows.Next() {
		var mrPickupTime MRPickupTime

		if err := rows.Scan(
			&mrPickupTime.Week,
			&mrPickupTime.Average,
			&mrPickupTime.Median,
			&mrPickupTime.Percentile75,
			&mrPickupTime.Percentile95,
		); err != nil {
			return nil, err
		}

		mrsPickupTime = append(mrsPickupTime, mrPickupTime)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mrsPickupTime, nil
}
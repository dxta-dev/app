package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type AggregatedCycleTimeStatistics struct {
	CodingTime Statistics                  `json:"coding_time"`
	PickupTime Statistics                  `json:"pickup_time"`
	ReviewTime Statistics                  `json:"review_time"`
	DeployTime Statistics                  `json:"deploy_time"`
	Weekly     []WeeklyCycleTimeStatistics `json:"weekly"`
}

type WeeklyCycleTimeStatistics struct {
	Week       string     `json:"week"`
	CodingTime Statistics `json:"coding_time"`
	PickupTime Statistics `json:"pickup_time"`
	ReviewTime Statistics `json:"review_time"`
	DeployTime Statistics `json:"deploy_time"`
}

type Statistics struct {
	Average      *float64 `json:"average"`
	Median       *float64 `json:"median"`
	Percentile75 *float64 `json:"percentile75"`
	Percentile95 *float64 `json:"percentile95"`
}

func DetailedCycleTime(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) (*AggregatedCycleTimeStatistics, error) {

	teamQuery := ""
	queryParamLength := len(weeks)

	if team != nil {
		teamQuery = "AND author.external_id in (SELECT member FROM tenant_team_members WHERE team = ?)"
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
		WITH dataset AS (
			SELECT
				mergedAt.week AS week,
				metrics.coding_duration AS coding_time,
				metrics.review_start_delay AS pickup_time,
				metrics.review_duration AS review_time,
				metrics.deploy_duration AS deploy_time
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
			AND metrics.deployed = 1
			AND repo.namespace_name = ?
			AND repo.name = ?
			AND branch.id = repo.default_branch
			%s
			AND author.bot = 0
		),
		data_by_week AS (
			SELECT
				week AS week,
				AVG(coding_time) AS avg_coding_time,
				AVG(pickup_time) AS avg_pickup_time,
				AVG(review_time) AS avg_review_time,
				AVG(deploy_time) AS avg_deploy_time,
				MEDIAN(coding_time) AS p50_coding_time,
				MEDIAN(pickup_time) AS p50_pickup_time,
				MEDIAN(review_time) AS p50_review_time,
				MEDIAN(deploy_time) AS p50_deploy_time,
				PERCENTILE_75(coding_time) AS p75_coding_time,
				PERCENTILE_75(pickup_time) AS p75_pickup_time,
				PERCENTILE_75(review_time) AS p75_review_time,
				PERCENTILE_75(deploy_time) AS p75_deploy_time,
				PERCENTILE_95(coding_time) AS p95_coding_time,
				PERCENTILE_95(pickup_time) AS p95_pickup_time,
				PERCENTILE_95(review_time) AS p95_review_time,
				PERCENTILE_95(deploy_time) AS p95_deploy_time
			FROM dataset
			GROUP BY week
		),
		data_total AS (
			SELECT
				AVG(coding_time) AS avg_coding_time,
				AVG(pickup_time) AS avg_pickup_time,
				AVG(review_time) AS avg_review_time,
				AVG(deploy_time) AS avg_deploy_time,
				MEDIAN(coding_time) AS p50_coding_time,
				MEDIAN(pickup_time) AS p50_pickup_time,
				MEDIAN(review_time) AS p50_review_time,
				MEDIAN(deploy_time) AS p50_deploy_time,
				PERCENTILE_75(coding_time) AS p75_coding_time,
				PERCENTILE_75(pickup_time) AS p75_pickup_time,
				PERCENTILE_75(review_time) AS p75_review_time,
				PERCENTILE_75(deploy_time) AS p75_deploy_time,
				PERCENTILE_95(coding_time) AS p95_coding_time,
				PERCENTILE_95(pickup_time) AS p95_pickup_time,
				PERCENTILE_95(review_time) AS p95_review_time,
				PERCENTILE_95(deploy_time) AS p95_deploy_time
			FROM dataset
		)
		SELECT
			NULL AS week,
			avg_coding_time,
			avg_pickup_time,
			avg_review_time,
			avg_deploy_time,
			p50_coding_time,
			p50_pickup_time,
			p50_review_time,
			p50_deploy_time,
			p75_coding_time,
			p75_pickup_time,
			p75_review_time,
			p75_deploy_time,
			p95_coding_time,
			p95_pickup_time,
			p95_review_time,
			p95_deploy_time
		FROM data_total
		UNION ALL
		SELECT
			week,
			avg_coding_time,
			avg_pickup_time,
			avg_review_time,
			avg_deploy_time,
			p50_coding_time,
			p50_pickup_time,
			p50_review_time,
			p50_deploy_time,
			p75_coding_time,
			p75_pickup_time,
			p75_review_time,
			p75_deploy_time,
			p95_coding_time,
			p95_pickup_time,
			p95_review_time,
			p95_deploy_time
		FROM data_by_week;
	`,
		weeksPlaceholder,
		teamQuery,
	)

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	wcts := make([]WeeklyCycleTimeStatistics, 0)
	acts := &AggregatedCycleTimeStatistics{}

	for rows.Next() {
		var week sql.NullString
		var coding_time Statistics
		var pickup_time Statistics
		var review_time Statistics
		var deploy_time Statistics
		if err := rows.Scan(
			&week,
			&coding_time.Average,
			&coding_time.Median,
			&coding_time.Percentile75,
			&coding_time.Percentile95,
			&pickup_time.Average,
			&pickup_time.Median,
			&pickup_time.Percentile75,
			&pickup_time.Percentile95,
			&review_time.Average,
			&review_time.Median,
			&review_time.Percentile75,
			&review_time.Percentile95,
			&deploy_time.Average,
			&deploy_time.Median,
			&deploy_time.Percentile75,
			&deploy_time.Percentile95,
		); err != nil {
			return nil, err
		}

		if week.Valid {
			wcts = append(wcts,
				WeeklyCycleTimeStatistics{
					CodingTime: coding_time,
					PickupTime: pickup_time,
					ReviewTime: review_time,
					DeployTime: deploy_time,
					Week:       week.String,
				},
			)
			continue
		}

		acts.CodingTime = coding_time
		acts.PickupTime = pickup_time
		acts.ReviewTime = review_time
		acts.DeployTime = deploy_time
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	acts.Weekly = wcts
	return acts, nil

}

func GetCycleTime(db *sql.DB, ctx context.Context, namespace string, repository string, weeks []string, team *int64) (*AggregatedStats, error) {

	teamQuery := ""
	queryParamLength := len(weeks)

	if team != nil {
		teamQuery = "AND author.external_id in (SELECT member FROM tenant_team_members WHERE team = ?)"
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

	query := buildQueryAggregatedStats(fmt.Sprintf(`
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
		AND metrics.deployed = 1
		AND repo.namespace_name = ?
		AND repo.name = ?
		AND branch.id = repo.default_branch
		%s
		AND author.bot = 0`,
		weeksPlaceholder,
		teamQuery,
	))

	rows, err := db.QueryContext(ctx, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cycleTimes, err := ScanAggregatedStatsRows(rows, weeks)

	if err != nil {
		return nil, err
	}

	return cycleTimes, nil
}
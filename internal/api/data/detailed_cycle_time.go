package data

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
)

type CycleTimeStatistics struct {
	CodingTime Statistics `json:"coding_time"`
	PickupTime Statistics `json:"pickup_time"`
	ReviewTime Statistics `json:"review_time"`
	DeployTime Statistics `json:"deploy_time"`
}

type AggregatedCycleTimeStatistics = OverallWeeklyData[CycleTimeStatistics]
type WeeklyCycleTimeStatistics = WeeklyData[CycleTimeStatistics]

func BuildDetailedCycleTimeQuery(weeks []string, team *int64) string {
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
			mergedAt.week AS week,
			metrics.coding_duration AS coding_time,
			metrics.review_start_delay AS pickup_time,
			metrics.review_duration AS review_time,
			CASE
				WHEN has_deployment.repository_external_id IS NULL THEN NULL
				WHEN metrics.deploy_duration = 0 THEN
					(unixepoch(date('now')) - unixepoch(
						CONCAT(dates.year, '-', LPAD(dates.month, 2, '0'), '-', LPAD(dates.day, 2, '0'))
					)) * 1000
				ELSE metrics.deploy_duration
			END AS deploy_time
    FROM transform_merge_request_metrics AS metrics
    JOIN transform_repositories AS repo
        ON repo.id = metrics.repository
		LEFT JOIN has_deployment
			ON has_deployment.repository_external_id = repo.external_id AND has_deployment.forge_type = repo.forge_type
    JOIN transform_merge_request_fact_dates_junk AS dj
        ON metrics.dates_junk = dj.id
    JOIN transform_dates AS mergedAt
        ON dj.merged_at = mergedAt.id
    JOIN transform_dates AS dates
        ON dj.merged_at = dates.id  -- Join the dates table to get the actual day, month, and year
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
}

func (d DB) GetDetailedCycleTime(ctx context.Context, query string, namespace string, repository string, weeks []string, team *int64) (*AggregatedCycleTimeStatistics, error) {
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

	rows, err := d.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wcts := make([]WeeklyCycleTimeStatistics, 0)
	acts := &AggregatedCycleTimeStatistics{}
	weekSet := make(map[string]bool)

	for rows.Next() {
		var week sql.NullString
		var codingTime, pickupTime, reviewTime, deployTime Statistics

		if err := rows.Scan(
			&week,
			&codingTime.Average,
			&pickupTime.Average,
			&reviewTime.Average,
			&deployTime.Average,
			&codingTime.Median,
			&pickupTime.Median,
			&reviewTime.Median,
			&deployTime.Median,
			&codingTime.Percentile75,
			&pickupTime.Percentile75,
			&reviewTime.Percentile75,
			&deployTime.Percentile75,
			&codingTime.Percentile95,
			&pickupTime.Percentile95,
			&reviewTime.Percentile95,
			&deployTime.Percentile95,
		); err != nil {
			return nil, err
		}

		if week.Valid {
			weekSet[week.String] = true
			wcts = append(wcts, WeeklyCycleTimeStatistics{
				Week: week.String,
				Data: CycleTimeStatistics{
					CodingTime: codingTime,
					PickupTime: pickupTime,
					ReviewTime: reviewTime,
					DeployTime: deployTime,
				},
			})
		} else {
			acts.Overall.CodingTime = codingTime
			acts.Overall.PickupTime = pickupTime
			acts.Overall.ReviewTime = reviewTime
			acts.Overall.DeployTime = deployTime
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, inputWeek := range weeks {
		if !weekSet[inputWeek] {
			wcts = append(wcts, WeeklyCycleTimeStatistics{
				Week: inputWeek,
				Data: CycleTimeStatistics{
					CodingTime: Statistics{},
					PickupTime: Statistics{},
					ReviewTime: Statistics{},
					DeployTime: Statistics{},
				},
			})
		}
	}

	sort.Slice(wcts, func(i, j int) bool {
		return wcts[i].Week < wcts[j].Week
	})

	acts.Weekly = wcts
	return acts, nil
}

package api

import (
	"database/sql"
	"fmt"
	"strings"
)

type MRSize struct {
	Week         string
	Average      int
	Median       int
	Percentile75 int
	Percentile95 int
	Count        int
}

/*
	SELECT
		FLOOR(AVG(metrics.mr_size)) as AVG,
		mergedAt.week as WEEK,
		COUNT(*) AS C,
		FLOOR(MEDIAN(metrics.mr_size)) AS P50,
		FLOOR(PERCENTILE_75(metrics.mr_size)) AS P75,
		FLOOR(PERCENTILE_95(metrics.mr_size)) as P95
	FROM transform_merge_request_metrics AS metrics
	JOIN transform_merge_request_fact_dates_junk AS dj
	ON metrics.dates_junk = dj.id
	JOIN transform_dates AS mergedAt
	ON dj.merged_at = mergedAt.id
	JOIN transform_merge_request_fact_users_junk AS uj
	ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
	ON uj.author = author.id
	WHERE mergedAt.week IN ("2024-W26", "2024-W27", "2024-W28", "2024-W29", "2024-W30", "2024-W31", "2024-W32", "2024-W33", "2024-W34", "2024-W35", "2024-W36", "2024-W37")
	AND author.external_id in (SELECT member FROM tenant_team_members WHERE team = 1)
	AND author.bot = 0
	GROUP BY mergedAt.week;
*/

func GetMRSize(db *sql.DB, weeks []string, team *int64) (map[string]MRSize, error) {
	return nil, nil
}

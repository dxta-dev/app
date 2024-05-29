package data

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"

	_ "github.com/libsql/libsql-client-go/libsql"
)

type AverageMRSizeByWeek struct {
	Week string
	Size int32
	N    int32
}

func NewAverageMRSizeByWeek(week string, size int32, n int32) AverageMRSizeByWeek {
	return AverageMRSizeByWeek{
		Week: week,
		Size: size,
		N:    n,
	}
}

func (s *Store) GetAverageMRSize(weeks []string, teamMembers []int64) (map[string]AverageMRSizeByWeek, float64, error) {
	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	query := fmt.Sprintf(`
	SELECT
		FLOOR(AVG(metrics.mr_size)),
		mergedAt.week,
		COUNT(*)
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
	AND author.bot = 0
	%s
	GROUP BY mergedAt.week;`,
		placeholders,
		usersInTeamConditionQuery)

	db, err := sql.Open(s.DriverName, s.DbUrl)

	if err != nil {
		return nil, 0, err
	}
	defer db.Close()

	queryParams := make([]interface{}, len(weeks)+len(teamMembers))
	for i, v := range weeks {
		queryParams[i] = v
	}
	for i, v := range teamMembers {
		queryParams[i+len(weeks)] = v
	}

	rows, err := db.QueryContext(s.Context, query, queryParams...)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	mrSizeByWeeks := make(map[string]AverageMRSizeByWeek)

	for rows.Next() {
		var avgMrSizeWeek AverageMRSizeByWeek

		if err := rows.Scan(&avgMrSizeWeek.Size, &avgMrSizeWeek.Week, &avgMrSizeWeek.N); err != nil {
			return nil, 0, err
		}

		mrSizeByWeeks[avgMrSizeWeek.Week] = NewAverageMRSizeByWeek(avgMrSizeWeek.Week, avgMrSizeWeek.Size, avgMrSizeWeek.N)
	}

	var totalMRSizeCount int32
	numOfWeeksWithMRSize := len(mrSizeByWeeks)

	for _, week := range weeks {
		totalMRSizeCount += mrSizeByWeeks[week].Size
		if _, ok := mrSizeByWeeks[week]; !ok {
			mrSizeByWeeks[week] = AverageMRSizeByWeek{
				Week: week,
				Size: 0,
				N:    0,
			}
		}
	}

	averageMRSizeByXWeeks := float64(totalMRSizeCount) / float64(numOfWeeksWithMRSize)

	return mrSizeByWeeks, averageMRSizeByXWeeks, nil
}

type AverageMrReviewDepthByWeek struct {
	Week  string
	Depth float32
}

func NewAverageMrReviewDepthByWeek(week string, depth float32) AverageMrReviewDepthByWeek {
	return AverageMrReviewDepthByWeek{
		Week:  week,
		Depth: depth,
	}
}

func (s *Store) GetAverageReviewDepth(weeks []string, teamMembers []int64) (map[string]AverageMrReviewDepthByWeek, float64, error) {
	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	query := fmt.Sprintf(`
	SELECT
		AVG(metrics.review_depth),
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
	WHERE mergedAt.week IN (%s)
	AND author.bot = 0
	%s
	GROUP BY mergedAt.week;`,
		placeholders,
		usersInTeamConditionQuery)

	db, err := sql.Open(s.DriverName, s.DbUrl)

	if err != nil {
		return nil, 0, err
	}

	defer db.Close()

	queryParams := make([]interface{}, len(weeks)+len(teamMembers))
	for i, v := range weeks {
		queryParams[i] = v
	}
	for i, v := range teamMembers {
		queryParams[i+len(weeks)] = v
	}

	rows, err := db.QueryContext(s.Context, query, queryParams...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	mrReviewDepthByWeeks := make(map[string]AverageMrReviewDepthByWeek)

	for rows.Next() {
		var avgReviewWeek AverageMrReviewDepthByWeek

		if err := rows.Scan(&avgReviewWeek.Depth, &avgReviewWeek.Week); err != nil {
			return nil, 0, err
		}

		mrReviewDepthByWeeks[avgReviewWeek.Week] = NewAverageMrReviewDepthByWeek(avgReviewWeek.Week, avgReviewWeek.Depth)
	}

	var totalReviewDepthCount float32
	numOfWeeksWithReviewDepth := len(mrReviewDepthByWeeks)

	for _, week := range weeks {
		totalReviewDepthCount += mrReviewDepthByWeeks[week].Depth
		if _, ok := mrReviewDepthByWeeks[week]; !ok {
			mrReviewDepthByWeeks[week] = AverageMrReviewDepthByWeek{
				Week:  week,
				Depth: 0,
			}
		}
	}

	averageReviewDepthByXWeeks := float64(totalReviewDepthCount) / float64(numOfWeeksWithReviewDepth)

	return mrReviewDepthByWeeks, averageReviewDepthByXWeeks, nil
}

type AverageHandoverPerMR struct {
	Week     string
	Handover float32
}

func NewAverageHandoverPerMR(week string, handover float32) AverageHandoverPerMR {
	return AverageHandoverPerMR{
		Week:     week,
		Handover: handover,
	}
}

func (s *Store) GetAverageHandoverPerMR(weeks []string, teamMembers []int64) (map[string]AverageHandoverPerMR, float64, error) {
	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	query := fmt.Sprintf(`
	SELECT
		AVG(metrics.handover),
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
	WHERE mergedAt.week IN (%s)
	AND author.bot = 0
	%s
	GROUP BY mergedAt.week;`,
		placeholders,
		usersInTeamConditionQuery)

	db, err := sql.Open(s.DriverName, s.DbUrl)

	if err != nil {
		return nil, 0, err
	}

	defer db.Close()

	queryParams := make([]interface{}, len(weeks)+len(teamMembers))
	for i, v := range weeks {
		queryParams[i] = v
	}
	for i, v := range teamMembers {
		queryParams[i+len(weeks)] = v
	}

	rows, err := db.QueryContext(s.Context, query, queryParams...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	mrHandoverByWeeks := make(map[string]AverageHandoverPerMR)

	for rows.Next() {
		var avgHandoverWeek AverageHandoverPerMR

		if err := rows.Scan(&avgHandoverWeek.Handover, &avgHandoverWeek.Week); err != nil {
			return nil, 0, err
		}

		mrHandoverByWeeks[avgHandoverWeek.Week] = NewAverageHandoverPerMR(avgHandoverWeek.Week, avgHandoverWeek.Handover)
	}

	var totalHandoverCount float32
	numOfWeeksWithHandover := len(mrHandoverByWeeks)

	for _, week := range weeks {
		totalHandoverCount += mrHandoverByWeeks[week].Handover
		if _, ok := mrHandoverByWeeks[week]; !ok {
			mrHandoverByWeeks[week] = AverageHandoverPerMR{
				Week:     week,
				Handover: 0,
			}
		}
	}

	averageHandoverByXWeeks := float64(totalHandoverCount) / float64(numOfWeeksWithHandover)

	return mrHandoverByWeeks, averageHandoverByXWeeks, nil
}

type MrCountByWeek struct {
	Week  string
	Count int32
}

func NewMrCountByWeek(week string, count int32) MrCountByWeek {
	return MrCountByWeek{
		Week:  week,
		Count: count,
	}
}

func (s *Store) GetMRsMergedWithoutReview(weeks []string, teamMembers []int64) (map[string]MrCountByWeek, float64, error) {
	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	query := fmt.Sprintf(`
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
	GROUP BY mergedAt.week;`,
		placeholders,
		usersInTeamConditionQuery)

	db, err := sql.Open(s.DriverName, s.DbUrl)

	if err != nil {
		return nil, 0, err
	}

	defer db.Close()

	queryParams := make([]interface{}, len(weeks)+len(teamMembers))
	for i, v := range weeks {
		queryParams[i] = v
	}
	for i, v := range teamMembers {
		queryParams[i+len(weeks)] = v
	}

	rows, err := db.QueryContext(s.Context, query, queryParams...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	mrCountByWeeks := make(map[string]MrCountByWeek)

	for rows.Next() {
		var mrCountWeek MrCountByWeek

		if err := rows.Scan(&mrCountWeek.Count, &mrCountWeek.Week); err != nil {
			return nil, 0, err
		}

		mrCountByWeeks[mrCountWeek.Week] = NewMrCountByWeek(mrCountWeek.Week, mrCountWeek.Count)
	}

	var totalMergedCount int32
	numOfWeeksWithMerged := len(mrCountByWeeks)

	for _, week := range weeks {
		totalMergedCount += mrCountByWeeks[week].Count
		if _, ok := mrCountByWeeks[week]; !ok {
			mrCountByWeeks[week] = MrCountByWeek{
				Week:  week,
				Count: 0,
			}
		}
	}

	averageMergedByXWeeks := float64(totalMergedCount) / float64(numOfWeeksWithMerged)

	return mrCountByWeeks, averageMergedByXWeeks, nil
}

func (s *Store) GetNewCodePercentage(weeks []string) (interface{}, error) {
	return nil, nil
}

func (s *Store) GetRefactorPercentage(weeks []string) (interface{}, error) {
	return nil, nil
}

func (s *Store) GetReworkPercentage(weeks []string) (interface{}, error) {
	return nil, nil
}

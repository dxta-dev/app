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
	Size int
	N    int
}

func (s *Store) GetAverageMRSize(weeks []string) (map[string]AverageMRSizeByWeek, float32, error) {

	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	query := fmt.Sprintf(`
	SELECT
		FLOOR(AVG(metrics.mr_size)),
		mergedAt.week,
		COUNT(*)
	FROM transform_merge_request_metrics as metrics
	JOIN transform_merge_request_fact_dates_junk as dj
	ON metrics.dates_junk = dj.id
	JOIN transform_dates as mergedAt
	ON dj.merged_at = mergedAt.id
	WHERE mergedAt.week IN (%s)
	GROUP BY mergedAt.week;`,
		placeholders)

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, 0, err
	}

	defer db.Close()

	weeksInterface := make([]interface{}, len(weeks))
	for i, v := range weeks {
		weeksInterface[i] = v
	}

	rows, err := db.Query(query, weeksInterface...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	mrSizeByWeeks := make(map[string]AverageMRSizeByWeek)

	for rows.Next() {
		var mrweek AverageMRSizeByWeek

		if err := rows.Scan(&mrweek.Size, &mrweek.Week, &mrweek.N); err != nil {
			return nil, 0, err
		}

		mrSizeByWeeks[mrweek.Week] = mrweek
	}

	totalMRSizeCount := 0

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

	averageMRSizeByXWeeks := float32(totalMRSizeCount) / float32(len(mrSizeByWeeks))

	return mrSizeByWeeks, averageMRSizeByXWeeks, nil
}

type AverageMrReviewDepthByWeek struct {
	Week  string
	Depth float32
}

func (s *Store) GetAverageReviewDepth(weeks []string) (map[string]AverageMrReviewDepthByWeek, float32, error) {
	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	query := fmt.Sprintf(`
	SELECT
		AVG(metrics.review_depth),
		mergedAt.week
	FROM transform_merge_request_metrics as metrics
	JOIN transform_merge_request_fact_dates_junk as dj
	ON metrics.dates_junk = dj.id
	JOIN transform_dates as mergedAt
	ON dj.merged_at = mergedAt.id
	WHERE mergedAt.week IN (%s)
	GROUP BY mergedAt.week;`,
		placeholders)

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, 0, err
	}

	defer db.Close()

	weeksInterface := make([]interface{}, len(weeks))
	for i, v := range weeks {
		weeksInterface[i] = v
	}

	rows, err := db.Query(query, weeksInterface...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	mrReviewDepthByWeeks := make(map[string]AverageMrReviewDepthByWeek)

	for rows.Next() {
		var mrweek AverageMrReviewDepthByWeek

		if err := rows.Scan(&mrweek.Depth, &mrweek.Week); err != nil {
			return nil, 0, err
		}

		mrReviewDepthByWeeks[mrweek.Week] = mrweek
	}

	totalReviewDepthCount := float32(0)

	for _, week := range weeks {
		totalReviewDepthCount += mrReviewDepthByWeeks[week].Depth
		if _, ok := mrReviewDepthByWeeks[week]; !ok {
			mrReviewDepthByWeeks[week] = AverageMrReviewDepthByWeek{
				Week:  week,
				Depth: 0,
			}
		}
	}

	averageReviewDepthByXWeeks := float32(totalReviewDepthCount) / float32(len(mrReviewDepthByWeeks))

	return mrReviewDepthByWeeks, averageReviewDepthByXWeeks, nil
}

type AverageHandoverPerMR struct {
	Week     string
	Handover float32
}

func (s *Store) GetAverageHandoverPerMR(weeks []string) (map[string]AverageHandoverPerMR, error) {
	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	query := fmt.Sprintf(`
	SELECT
		AVG(metrics.handover),
		mergedAt.week
	FROM transform_merge_request_metrics as metrics
	JOIN transform_merge_request_fact_dates_junk as dj
	ON metrics.dates_junk = dj.id
	JOIN transform_dates as mergedAt
	ON dj.merged_at = mergedAt.id
	WHERE mergedAt.week IN (%s)
	GROUP BY mergedAt.week;`,
		placeholders)

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	weeksInterface := make([]interface{}, len(weeks))
	for i, v := range weeks {
		weeksInterface[i] = v
	}

	rows, err := db.Query(query, weeksInterface...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mrHandoverByWeeks := make(map[string]AverageHandoverPerMR)

	for rows.Next() {
		var mrweek AverageHandoverPerMR

		if err := rows.Scan(&mrweek.Handover, &mrweek.Week); err != nil {
			return nil, err
		}

		mrHandoverByWeeks[mrweek.Week] = mrweek
	}

	for _, week := range weeks {
		if _, ok := mrHandoverByWeeks[week]; !ok {
			mrHandoverByWeeks[week] = AverageHandoverPerMR{
				Week:     week,
				Handover: 0,
			}
		}
	}

	return mrHandoverByWeeks, nil
}

type MrCountByWeek struct {
	Week  string
	Count int
}

func (s *Store) GetMRsMergedWithoutReview(weeks []string) (map[string]MrCountByWeek, float32, error) {
	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	query := fmt.Sprintf(`
	SELECT
		COUNT(*),
		mergedAt.week
	FROM transform_merge_request_metrics as metrics
	JOIN transform_merge_request_fact_dates_junk as dj
	ON metrics.dates_junk = dj.id
	JOIN transform_dates as mergedAt
	ON dj.merged_at = mergedAt.id
	WHERE mergedAt.week IN (%s) and metrics.review_depth = 0
	GROUP BY mergedAt.week;`,
		placeholders)

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, 0, err
	}

	defer db.Close()

	weeksInterface := make([]interface{}, len(weeks))
	for i, v := range weeks {
		weeksInterface[i] = v
	}

	rows, err := db.Query(query, weeksInterface...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	mrCountByWeeks := make(map[string]MrCountByWeek)

	for rows.Next() {
		var mrweek MrCountByWeek

		if err := rows.Scan(&mrweek.Count, &mrweek.Week); err != nil {
			return nil, 0, err
		}

		mrCountByWeeks[mrweek.Week] = mrweek
	}

	totalMergedCount := 0

	for _, week := range weeks {
		totalMergedCount += mrCountByWeeks[week].Count
		if _, ok := mrCountByWeeks[week]; !ok {
			mrCountByWeeks[week] = MrCountByWeek{
				Week:  week,
				Count: 0,
			}
		}
	}

	averageMergedByXWeeks := float32(totalMergedCount) / float32(len(mrCountByWeeks))

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

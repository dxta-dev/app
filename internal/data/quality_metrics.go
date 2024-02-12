package data

import (
	"fmt"
	"strings"

	"database/sql"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type AverageMRSizeByWeek struct {
	Week string
	Size int
	N    int
}

func (s *Store) GetAverageMRSize(weeks []string) ([]AverageMRSizeByWeek, error) {

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

	var mrSizeByWeeks []AverageMRSizeByWeek

	for rows.Next() {
		var mrweek AverageMRSizeByWeek

		if err := rows.Scan(
			&mrweek.Size,
			&mrweek.Week,
			&mrweek.N,
		); err != nil {
			return nil, err
		}

		mrSizeByWeeks = append(mrSizeByWeeks, mrweek)
	}

	return mrSizeByWeeks, nil
}

type AverageMrReviewDepthByWeek struct {
	Week  string
	Depth float32
}

func (s *Store) GetAverageReviewDepth(weeks []string) ([]AverageMrReviewDepthByWeek, error) {
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

	var mrReviewDepthByWeeks []AverageMrReviewDepthByWeek

	for rows.Next() {
		var mrweek AverageMrReviewDepthByWeek

		if err := rows.Scan(
			&mrweek.Depth,
			&mrweek.Week,
		); err != nil {
			return nil, err
		}

		mrReviewDepthByWeeks = append(mrReviewDepthByWeeks, mrweek)
	}

	return mrReviewDepthByWeeks, nil
}

type MrCountByWeek struct {
	Week  string
	Count int
}

func (s *Store) GetMRsMergedWithoutReview(weeks []string) ([]MrCountByWeek, error) {
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

	var mrCountByWeeks []MrCountByWeek

	for rows.Next() {
		var mrweek MrCountByWeek

		if err := rows.Scan(
			&mrweek.Count,
			&mrweek.Week,
		); err != nil {
			return nil, err
		}

		mrCountByWeeks = append(mrCountByWeeks, mrweek)
	}

	return mrCountByWeeks, nil
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

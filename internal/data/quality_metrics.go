package data

import (
	"fmt"
	"strings"

	"database/sql"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type PRSizeByWeek struct {
	Week string
	Size int
}

func (s *Store) GetAverageMRSize(weeks []string) ([]PRSizeByWeek, error) {

	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	query := fmt.Sprintf(`
	SELECT
		FLOOR(AVG(metrics.mr_size)),
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

	var prSizeByWeeks []PRSizeByWeek

	for rows.Next() {
		var prweek PRSizeByWeek

		if err := rows.Scan(
			&prweek.Size,
			&prweek.Week,
		); err != nil {
			return nil, err
		}

		prSizeByWeeks = append(prSizeByWeeks, prweek)
	}

	return prSizeByWeeks, nil
}

func (s *Store) GetAverageReviewDepth(weeks []string) (interface{}, error) {
	return nil, nil
}

func (s *Store) GetPRsMergedWithoutReview(weeks []string) (interface{}, error) {
	return nil, nil
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

package data

import (
	"database/sql"
)

type MergeRequestMetrics struct {
	MrSize      int
	Handover    int
	ReviewDepth int
}

func (s *Store) GetMergeRequestMetrics(mrId int64) (*MergeRequestMetrics, error) {
	query := `
        SELECT metrics.mr_size, metrics.handover, metrics.review_depth FROM transform_merge_request_metrics AS metrics
        WHERE metrics.merge_request = ?
    `

	db, err := sql.Open("libsql", s.DbUrl)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var mrMetrics MergeRequestMetrics
	row := db.QueryRow(query, mrId)
	err = row.Scan(&mrMetrics.MrSize, &mrMetrics.Handover, &mrMetrics.ReviewDepth)
	if err != nil {
		return nil, err
	}

	return &mrMetrics, nil
}

func (s *Store) GetTotalCommitsForMR(mrId int64) (int, error) {
	query := `
	SELECT COUNT (*) FROM transform_merge_request_events AS events
	WHERE events.merge_request_event_type = 9
	AND events.merge_request = ?
`

	db, err := sql.Open("libsql", s.DbUrl)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var totalCommits int
	row := db.QueryRow(query, mrId)
	err = row.Scan(&totalCommits)
	if err != nil {
		return 0, err
	}

	return totalCommits, nil
}

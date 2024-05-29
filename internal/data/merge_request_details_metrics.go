package data

import (
	"database/sql"
)

type MergeRequestMetrics struct {
	MrSize      int
	Handover    int
	ReviewDepth int
}

type MergeRequestData struct {
	Metrics      MergeRequestMetrics
	TotalCommits int
}

func (s *Store) GetMergeRequestMetricsData(mrId int64) (*MergeRequestData, error) {
	metricsQuery := `
        SELECT metrics.mr_size, metrics.handover, metrics.review_depth FROM transform_merge_request_metrics AS metrics
        WHERE metrics.merge_request = ?
    `
	commitsQuery := `
        SELECT COUNT(*) FROM transform_merge_request_events AS events
        WHERE events.merge_request_event_type = 9
        AND events.merge_request = ?
    `

	db, err := sql.Open(s.DriverName, s.DbUrl)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var mrData MergeRequestData

	metricsRow := db.QueryRowContext(s.Context, metricsQuery, mrId)
	err = metricsRow.Scan(&mrData.Metrics.MrSize, &mrData.Metrics.Handover, &mrData.Metrics.ReviewDepth)
	if err != nil {
		return nil, err
	}

	commitsRow := db.QueryRowContext(s.Context, commitsQuery, mrId)
	err = commitsRow.Scan(&mrData.TotalCommits)
	if err != nil {
		return nil, err
	}

	return &mrData, nil
}

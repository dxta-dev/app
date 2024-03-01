package data

import (
	"database/sql"
	"fmt"
	"strings"
)

type CodeChangesCount struct {
	Count int
	Week  string
}

func (s *Store) GetTotalCodeChanges(weeks []string) (map[string]CodeChangesCount, float32, error) {

	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	query := fmt.Sprintf(`
	SELECT
		SUM(metrics.mr_size) AS total_mr_size,
		dates.week
	FROM transform_merge_request_metrics as metrics
	JOIN transform_merge_request_fact_dates_junk as dates_junk
	ON metrics.dates_junk = dates_junk.id
	JOIN transform_dates as dates
	ON dates_junk.merged_at = dates.id
	WHERE dates.week IN (%s)
	GROUP BY dates.week;`,
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

	codeChangesByWeek := make(map[string]CodeChangesCount)

	for rows.Next() {
		var codeChangesCount CodeChangesCount

		if err = rows.Scan(&codeChangesCount.Count, &codeChangesCount.Week); err != nil {
			return nil, 0, err
		}
		codeChangesByWeek[codeChangesCount.Week] = codeChangesCount
	}

	totalCodeChangesCount := 0

	for _, week := range weeks {
		totalCodeChangesCount += codeChangesByWeek[week].Count
		if _, ok := codeChangesByWeek[week]; !ok {
			codeChangesByWeek[week] = CodeChangesCount{
				Count: 0,
				Week:  week,
			}
		}
	}

	averageCodeChangesByXWeeks := float32(totalCodeChangesCount) / float32(len(weeks))

	return codeChangesByWeek, averageCodeChangesByXWeeks, nil
}

type CommitCountByWeek struct {
	Week  string
	Count int
}

func (s *Store) GetTotalCommits(weeks []string) (map[string]CommitCountByWeek, error) {

	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	query := fmt.Sprintf(`
		SELECT
			COUNT (*),
			commitedAt.week
		FROM transform_merge_request_events as ev
		JOIN transform_dates as commitedAt
		ON commitedAt.id = ev.commited_at
		WHERE ev.merge_request_event_type = 9
		AND commitedAt.week IN (%s)
		GROUP BY commitedAt.week;`,
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

	commitCountByWeeks := make(map[string]CommitCountByWeek)

	for rows.Next() {
		var commitCount CommitCountByWeek

		if err := rows.Scan(&commitCount.Count, &commitCount.Week); err != nil {
			return nil, err
		}
		commitCountByWeeks[commitCount.Week] = commitCount
	}

	for _, week := range weeks {
		if _, ok := commitCountByWeeks[week]; !ok {
			commitCountByWeeks[week] = CommitCountByWeek{
				Week:  week,
				Count: 0,
			}
		}
	}

	return commitCountByWeeks, nil
}

func (s *Store) GetTotalMrsOpened(weeks []string) (map[string]MrCountByWeek, error) {

	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	query := fmt.Sprintf(`
	SELECT
		COUNT (*),
		opened_dates.week
	FROM transform_merge_request_metrics as metrics
	JOIN transform_merge_request_fact_dates_junk as dates_junk
	ON metrics.dates_junk = dates_junk.id
	JOIN transform_dates as opened_dates
	ON dates_junk.opened_at = opened_dates.id
	WHERE opened_dates.week IN (%s)
	GROUP BY opened_dates.week`,
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

	prCountByWeeks := make(map[string]MrCountByWeek)

	for rows.Next() {
		var prCount MrCountByWeek

		if err := rows.Scan(&prCount.Count, &prCount.Week); err != nil {
			return nil, err
		}
		prCountByWeeks[prCount.Week] = prCount
	}

	for _, week := range weeks {
		if _, ok := prCountByWeeks[week]; !ok {
			prCountByWeeks[week] = MrCountByWeek{
				Week:  week,
				Count: 0,
			}
		}
	}

	return prCountByWeeks, nil
}

type TotalReviewsByWeek struct {
	Week  string
	Count int
}

func (s *Store) GetTotalReviews(weeks []string) (map[string]TotalReviewsByWeek, error) {

	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	query := fmt.Sprintf(`
		SELECT
			COUNT (*),
			occuredAt.week
		FROM transform_merge_request_events as ev
		JOIN transform_dates as occuredAt
		ON occuredAt.id = ev.occured_on
		WHERE ev.merge_request_event_type = 15
		AND occuredAt.week IN (%s)
		GROUP BY occuredAt.week;`,
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

	totalReviewsByWeek := make(map[string]TotalReviewsByWeek)

	for rows.Next() {
		var reviewCount TotalReviewsByWeek

		if err := rows.Scan(&reviewCount.Count, &reviewCount.Week); err != nil {
			return nil, err
		}
		totalReviewsByWeek[reviewCount.Week] = reviewCount
	}

	for _, week := range weeks {
		if _, ok := totalReviewsByWeek[week]; !ok {
			totalReviewsByWeek[week] = TotalReviewsByWeek{
				Week:  week,
				Count: 0,
			}
		}
	}

	return totalReviewsByWeek, nil
}

type MergeFrequencyByWeek struct {
	Week   string
	Amount float32
}

func (s *Store) GetMergeFrequency(weeks []string) (map[string]MergeFrequencyByWeek, error) {
	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	query := fmt.Sprintf(`
		SELECT
			CAST(COUNT (*) AS REAL) / 7,
			merged_dates.week
		FROM transform_merge_request_metrics as metrics
		JOIN transform_merge_request_fact_dates_junk as dates_junk
		ON metrics.dates_junk = dates_junk.id
		JOIN transform_dates as merged_dates
		ON dates_junk.merged_at = merged_dates.id
		WHERE merged_dates.week IN (%s)
		GROUP BY merged_dates.week`,
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

	mergeFrequencyByWeek := make(map[string]MergeFrequencyByWeek)

	for rows.Next() {
		var mergeFreq MergeFrequencyByWeek

		if err := rows.Scan(&mergeFreq.Amount, &mergeFreq.Week); err != nil {
			return nil, err
		}
		mergeFrequencyByWeek[mergeFreq.Week] = mergeFreq
	}

	for _, week := range weeks {
		if _, ok := mergeFrequencyByWeek[week]; !ok {
			mergeFrequencyByWeek[week] = MergeFrequencyByWeek{
				Week:   week,
				Amount: 0,
			}
		}
	}

	return mergeFrequencyByWeek, nil
}

func (s *Store) GetDeployFrequency(weeks []string) (interface{}, error) {
	return nil, nil
}

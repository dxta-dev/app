package data

import (
	"database/sql"
	"fmt"
	"strings"
)

func (s *Store) GetTotalCodeChanges(weeks []string) (interface{}, error) {
	return nil, nil
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
		ON commitedAt.id = ev.occured_on
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

func (s *Store) GetTotalPRsOpened(weeks []string) (interface{}, error) {
	return nil, nil
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

func (s *Store) GetMergeFrequency(weeks []string) (interface{}, error) {
	return nil, nil
}

func (s *Store) GetDeployFrequency(weeks []string) (interface{}, error) {
	return nil, nil
}

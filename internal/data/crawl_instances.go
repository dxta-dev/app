package data

import (
	"database/sql"
	"sort"
	"time"
)

type TimeFrame struct {
	Since time.Time
	Until time.Time
}

type TimeFrameSlice []TimeFrame

func (tfs TimeFrameSlice) Len() int {
	return len(tfs)
}

func (tfs TimeFrameSlice) Less(i, j int) bool {
	return tfs[i].Since.Before(tfs[j].Since)
}

func (tfs TimeFrameSlice) Swap(i, j int) {
	tfs[i], tfs[j] = tfs[j], tfs[i]
}

func (s *Store) GetCrawlInstances(from, to time.Time) (TimeFrameSlice, error) {

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	query := `
		SELECT since, until
		FROM crawl_instances
		WHERE since >= ?
		AND until <= ?;
		`

	rows, err := db.Query(query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var crawlInstances TimeFrameSlice

	for rows.Next() {
		var since, until time.Time
		if err := rows.Scan(&since, &until); err != nil {
			return nil, err
		}
		crawlInstances = append(crawlInstances, TimeFrame{Since: since, Until: until})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return crawlInstances, nil
}

func FindGaps(from, to time.Time, timeFrames TimeFrameSlice) TimeFrameSlice {
	var gaps TimeFrameSlice

	sort.Sort(timeFrames)

	var until time.Time
	if len(timeFrames) == 0 {
		until = to
	} else {
		until = timeFrames[0].Since
	}

	gaps = append(gaps, TimeFrame{
		Since: from,
		Until: until,
	})

	for i := 1; i < len(timeFrames)-1; i++ {
		if timeFrames[i-1].Until.Before(timeFrames[i].Since) {
			gaps = append(gaps, TimeFrame{
				Since: timeFrames[i-1].Until,
				Until: timeFrames[i].Since,
			})
		}
	}

	return gaps
}

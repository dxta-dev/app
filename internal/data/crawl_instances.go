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

func (s *Store) GetCrawlInstances(from, to int64) (TimeFrameSlice, error) {

	db, err := sql.Open("libsql", s.DbUrl)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
        SELECT since, until
        FROM crawl_instances
        WHERE
        (
            (since < ? AND until > ?)
            OR
            (since > ? AND until > ? AND until < ?)
            OR
            (since > ? AND since < ? AND until > ? AND until < ?)
            OR
            (since < ? AND since > ? AND until > ?)
        )
        `

	rows, err := db.Query(query, from, to, from, from, to, from, to, from, to, to, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var crawlInstances TimeFrameSlice

	for rows.Next() {
		var sinceInt, untilInt int64
		if err := rows.Scan(&sinceInt, &untilInt); err != nil {
			return nil, err
		}

		since := time.Unix(sinceInt, 0)
		until := time.Unix(untilInt, 0)

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

	currentFrame := TimeFrame{Since: from, Until: from}

	for _, frame := range timeFrames {

		if currentFrame.Until.Before(frame.Since) {
			gaps = append(gaps, TimeFrame{Since: currentFrame.Until, Until: frame.Since})
		}

		if frame.Until.After(currentFrame.Until) {
			currentFrame.Until = frame.Until
		}
	}

	if currentFrame.Until.Before(to) {
		gaps = append(gaps, TimeFrame{Since: currentFrame.Until, Until: to})
	}

	return gaps
}

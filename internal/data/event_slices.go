package data

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

const (
	UNKNOWN EventType = iota
	OPENED
	STARTED_CODING
	STARTED_PICKUP
	STARTED_REVIEW
	NOTED
	ASSIGNED
	CLOSED
	COMMENTED
	COMMITTED
	CONVERT_TO_DRAFT
	MERGED
	READY_FOR_REVIEW
	REVIEW_REQUEST_REMOVED
	REVIEW_REQUESTED
	REVIEWED
	UNASSIGNED
)

type EventType int

type Event struct {
	Timestamp       int64
	Type            EventType
	Actor           string
	MergeRequest    string
	MergeRequestUrl string
}

type EventSlice []Event

func (d EventSlice) Len() int {
	return len(d)
}

func (d EventSlice) Less(i, j int) bool {
	return d[i].Timestamp < d[j].Timestamp
}

func (d EventSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (s *Store) GetEventSlices(date time.Time) (EventSlice, error) {
	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, err
	}

	year, week := date.ISOWeek()

	query := `
		SELECT
			u.name,
			mr.title,
			mr.web_url,
			ev.timestamp,
			ev.merge_request_event_type
		FROM transform_merge_request_events as ev
		JOIN transform_dates as d ON d.id = ev.occured_on
		JOIN transform_forge_users as u ON u.id = ev.actor
		JOIN transform_merge_requests as mr ON mr.id = ev.merge_request
		WHERE d.week=? AND d.year=?;
	`
	rows, err := db.Query(query, week, year)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event

		var timestamp int64

		var eventType int

		var actor string

		var mergeRequest string

		var mergeRequestUrl string

		if err := rows.Scan(&actor, &mergeRequest, &mergeRequestUrl, &timestamp, &eventType); err != nil {
			log.Fatal(err)
		}

		event.Type = EventType(eventType)
		event.Timestamp = timestamp
		event.Actor = actor
		event.MergeRequest = mergeRequest
		event.MergeRequestUrl = mergeRequestUrl

		events = append(events, event)
	}

	return events, nil
}

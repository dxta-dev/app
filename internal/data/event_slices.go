package data

import (
	"database/sql"
	"dxta-dev/app/internal/utils"
	"fmt"
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
	Id              int64
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

func (s *Store) GetEventInfo(eventId int64) (*Event, error) {
	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, err
	}

	query := `SELECT
			ev.id,
			user.name,
			mr.title,
			mr.web_url,
			ev.timestamp,
			ev.merge_request_event_type
		FROM transform_merge_request_events as ev
		JOIN transform_dates as date ON date.id = ev.occured_on
		JOIN transform_forge_users as user ON user.id = ev.actor
		JOIN transform_merge_requests as mr ON mr.id = ev.merge_request
		WHERE ev.id = ?;
		`
	rows, err := db.Query(query, eventId)

	if err != nil {
		return nil, err
	}
	fmt.Println("rovoviiii", rows)

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event

		if err := rows.Scan(&event.Id, &event.Actor, &event.MergeRequest, &event.MergeRequestUrl, &event.Timestamp, &event.Type); err != nil {
			log.Fatal(err)
		}

		events = append(events, event)
	}

	event := &Event{
		Id:              1,
		Timestamp:       123456789,
		Type:            OPENED,
		Actor:           "John Doe",
		MergeRequest:    "123",
		MergeRequestUrl: "http://example.com",
	}

	return event, nil

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

	week := utils.GetFormattedWeek(date)

	query := `
		SELECT
			ev.id,
			user.name,
			mr.title,
			mr.web_url,
			ev.timestamp,
			ev.merge_request_event_type
		FROM transform_merge_request_events as ev
		JOIN transform_dates as date ON date.id = ev.occured_on
		JOIN transform_forge_users as user ON user.id = ev.actor
		JOIN transform_merge_requests as mr ON mr.id = ev.merge_request
		WHERE date.week=?;
	`
	rows, err := db.Query(query, week)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event

		if err := rows.Scan(&event.Id, &event.Actor, &event.MergeRequest, &event.MergeRequestUrl, &event.Timestamp, &event.Type); err != nil {
			log.Fatal(err)
		}

		events = append(events, event)
	}

	return events, nil
}

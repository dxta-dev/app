package data

import (
	"database/sql"
	"dxta-dev/app/internal/utils"
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
	Id                int64
	Timestamp         int64
	Type              EventType
	Actor             int64
	MergeRequestId    int64
	MergeRequestTitle string
	MergeRequestUrl   string
}

type EventSlice []Event

func (d EventSlice) Len() int {
	return len(d)
}

func (d EventSlice) Less(i, j int) bool {
	return d[i].Timestamp < d[j].Timestamp || (d[i].Timestamp == d[j].Timestamp && d[i].Type < d[j].Type)
}

func (d EventSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (s *Store) GetMergeRequestEvents(mrId int64) (EventSlice, error) {
	return nil, nil
}

func (s *Store) GetEventSlices(date time.Time) (EventSlice, error) {
	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	week := utils.GetFormattedWeek(date)

	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			ev.id,
			user.id,
			mr.id,
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

		if err := rows.Scan(
			&event.Id, &event.Actor, &event.MergeRequestId,
			&event.MergeRequestTitle, &event.MergeRequestUrl,
			&event.Timestamp, &event.Type,
		); err != nil {
			log.Fatal(err)
		}

		events = append(events, event)
	}

	return events, nil
}

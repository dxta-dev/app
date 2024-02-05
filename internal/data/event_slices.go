package data

import (
	"database/sql"
	"dxta-dev/app/internal/utils"
	"log"
	"sort"
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

func (d EventSlice) GroupByMergeRequest() []EventSlice {
	groupMap := make(map[int64]EventSlice)
	var grouped []EventSlice

	for _, event := range d {
		groupMap[event.MergeRequestId] = append(groupMap[event.MergeRequestId], event)
	}

	var keys []int64

	for k := range groupMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, key := range keys {
		sort.Slice(groupMap[key], func(i, j int) bool {
			return groupMap[key][i].Timestamp < groupMap[key][j].Timestamp
		})
		grouped = append(grouped, groupMap[key])
	}

	return grouped
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

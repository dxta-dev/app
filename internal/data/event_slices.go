package data

import (
	"fmt"

	"github.com/dxta-dev/app/internal/util"

	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"

	_ "github.com/libsql/libsql-client-go/libsql"
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

type EventUserInfo struct {
	Id         int64
	Name       string
	ProfileUrl string
	AvatarUrl  string
}

type Event struct {
	Id                  int64
	Timestamp           int64
	Type                EventType
	Actor               EventUserInfo
	MergeRequestId      int64
	MergeRequestCanonId int64
	MergeRequestTitle   string
	MergeRequestUrl     string
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
	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	query := `
		SELECT
			ev.id,
			user.id,
			user.profile_url,
			user.avatar_url,
			user.name,
			mr.id,
			mr.canon_id,
			mr.title,
			mr.web_url,
			ev.timestamp,
			ev.merge_request_event_type
		FROM transform_merge_request_events AS ev
		JOIN transform_forge_users AS user ON user.id = ev.actor
		JOIN transform_merge_requests AS mr ON mr.id = ev.merge_request
		WHERE ev.merge_request =?
		AND user.bot = 0
		ORDER BY ev.timestamp ASC;
		`

	rows, err := db.Query(query, mrId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mergeRequestEvents EventSlice

	for rows.Next() {
		var event Event

		if err := rows.Scan(
			&event.Id,
			&event.Actor.Id, &event.Actor.ProfileUrl, &event.Actor.AvatarUrl, &event.Actor.Name,
			&event.MergeRequestId, &event.MergeRequestCanonId, &event.MergeRequestTitle, &event.MergeRequestUrl,
			&event.Timestamp, &event.Type,
		); err != nil {
			log.Fatal(err)
		}

		mergeRequestEvents = append(mergeRequestEvents, event)
	}

	return mergeRequestEvents, nil
}

func (s *Store) GetEventSlices(date time.Time, team *TeamRef) (EventSlice, error) {
	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	week := util.GetFormattedWeek(date)

	query := fmt.Sprintf(`
		SELECT
		ev.id,
		user.id,
		mr.id,
		mr.title,
		mr.web_url,
		ev.timestamp,
		ev.merge_request_event_type
	FROM transform_merge_request_events AS ev
	JOIN transform_dates AS date ON date.id = ev.occured_on
	JOIN transform_forge_users AS user ON user.id = ev.actor
	JOIN transform_merge_requests AS mr ON mr.id = ev.merge_request
	JOIN transform_merge_request_metrics AS metrics ON metrics.merge_request = mr.id
	JOIN transform_merge_request_fact_users_junk AS u ON u.id = metrics.users_junk
	JOIN transform_forge_users AS author ON author.id = u.author
	WHERE date.week = '%v'
	AND author.bot = 0
	AND user.bot = 0%s;
		`, week, AndUserInTeamQueryPart("user.external_id", team))

	rows, err := db.Query(query, week)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event

		if err := rows.Scan(
			&event.Id, &event.Actor.Id, &event.MergeRequestId,
			&event.MergeRequestTitle, &event.MergeRequestUrl,
			&event.Timestamp, &event.Type,
		); err != nil {
			log.Fatal(err)
		}

		events = append(events, event)
	}

	return events, nil
}

package handlers

import (
	"database/sql"
	"dxta-dev/app/internal/graphs"
	"dxta-dev/app/internal/templates"
	"dxta-dev/app/internal/utils"
	"log"
	"sort"
	"time"

	"github.com/joho/godotenv"
	"github.com/wcharczuk/go-chart/v2/drawing"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

const (
	UNKNOWN templates.EventType = iota
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

type EventSlice []templates.Event

func (d EventSlice) Len() int {
	return len(d)
}

func (d EventSlice) Less(i, j int) bool {
	return d[i].Timestamp < d[j].Timestamp
}

func (d EventSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func getData(date time.Time, dbUrl string) (EventSlice, error) {

	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	db, err := sql.Open("libsql", dbUrl)

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
			ev.actor,
			ev.timestamp,
			ev.merge_request_event_type
		FROM transform_merge_request_events as ev
		JOIN transform_dates as d ON d.id = ev.occured_on
		WHERE d.week=? AND d.year=?;
	`
	rows, err := db.Query(query, week, year)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []templates.Event

	for rows.Next() {
		var event templates.Event

		var timestamp int64

		var eventType int

		var actor int64

		if err := rows.Scan(&actor, &timestamp, &eventType); err != nil {
			log.Fatal(err)
		}

		event.Type = templates.EventType(eventType)
		event.Timestamp = timestamp
		event.Actor = actor

		events = append(events, event)
	}

	return events, nil
}

func getSeries(date time.Time, dbUrl string) templates.SwarmSeries {
	var xvalues []float64
	var yvalues []float64

	startOfWeek := utils.GetStartOfWeek(date)

	var times []time.Time

	events, _ := getData(date, dbUrl)

	sort.Sort(events)

	for _, d := range events {
		t := time.Unix(d.Timestamp/1000, 0)
		times = append(times, t)
	}

	for _, t := range times {
		xSecondsValue := float64(t.Unix() - startOfWeek.Unix())
		xvalues = append(xvalues, xSecondsValue)
		yvalues = append(yvalues, 60*60*12)
	}

	xvalues, yvalues = graphs.Beehive(xvalues, yvalues, 1400, 200, 5)

	colors := []drawing.Color{}

	for i := 0; i < len(xvalues); i++ {
		switch events[i].Type {
		case COMMITTED:
			colors = append(colors, drawing.ColorBlue)
		case MERGED:
			colors = append(colors, drawing.ColorRed)
		case REVIEWED:
			colors = append(colors, drawing.ColorGreen)
		default:
			colors = append(colors, drawing.ColorFromAlphaMixedRGBA(204, 204, 204, 255))
		}
	}

	return templates.SwarmSeries{
		Title:     "series 1",
		DotColors: colors,
		XValues:   xvalues,
		YValues:   yvalues,
	}
}

package handlers

import (
	"context"
	"database/sql"
	"dxta-dev/app/internals/graphs"
	"dxta-dev/app/internals/templates"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/wcharczuk/go-chart/v2/drawing"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type EventType int

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

type Event struct {
	Timestamp int64
	Type      EventType
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

func getData(date time.Time) (EventSlice, error) {

	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	db, err := sql.Open("libsql", os.Getenv("DATABASE_URL"))

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

	var events []Event

	for rows.Next() {
		var event Event

		var timestamp int64

		var eventType int

		if err := rows.Scan(&timestamp, &eventType); err != nil {
			log.Fatal(err)
		}

		event.Type = EventType(eventType)
		event.Timestamp = timestamp
		events = append(events, event)
	}

	return events, nil
}

func getSeries(date time.Time) templates.SwarmSeries {
	var xvalues []float64
	var yvalues []float64

	startOfWeek := GetStartOfWeek(date)

	var times []time.Time

	events, _ := getData(time.Now())

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

func GetStartOfWeek(date time.Time) time.Time {
	offset := int(time.Monday - date.Weekday())

	if offset > 0 {
		offset = -6
	}

	startOfWeek := date.AddDate(0, 0, offset)

	startOfWeek = startOfWeek.Truncate(24 * time.Hour)

	return startOfWeek
}

func GetCurrentWeek(date time.Time) string {
	year, week := date.ISOWeek()

	formattedWeek := fmt.Sprintf("%d-W%02d", year, week)

	return formattedWeek
}

func parseYearWeek(yw string) (time.Time, error) {
	parts := strings.Split(yw, "-W")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid format")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, err
	}

	week, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, err
	}

	firstDayOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	offset := int(time.Monday - firstDayOfYear.Weekday())
	if offset > 0 {
		offset -= 7
	}

	daysToStartOfWeek := (week-1)*7 + offset
	startOfWeek := firstDayOfYear.AddDate(0, 0, daysToStartOfWeek)

	return startOfWeek, nil
}

func (a *App) Swarm(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	page := &templates.Page{
		Title:   "Charts",
		Boosted: h.HxBoosted,
	}

	date := time.Now()

	weekString := r.URL.Query().Get("week")

	fmt.Println("weekString", weekString)

	if weekString != "" {
		dateTime, err := parseYearWeek(weekString)

		if err == nil {
			date = dateTime
		} else {
			res := c.Response()

			res.Header().Set("HX-Push-Url", "/swarm?week="+weekString)
		}
	}

	fmt.Println(date)
	startOfWeek := GetStartOfWeek(date)
	fmt.Println(startOfWeek)


	if h.HxRequest && h.HxTrigger != "" {
		components := templates.SwarmChart(getSeries(date), startOfWeek)
		return components.Render(context.Background(), c.Response().Writer)
	}

	components := templates.Swarm(page, getSeries(date), startOfWeek, GetCurrentWeek(date))

	return components.Render(context.Background(), c.Response().Writer)
}

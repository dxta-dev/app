package handlers

import (
	"context"
	"database/sql"
	"dxta-dev/app/internals/templates"
	"fmt"
	"os"
	"strings"

	"github.com/donseba/go-htmx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type JoinedIndexes struct {
	Id              int
	UserJunk        int
	MrSize          int
	CodingDuration  int
	PickupDuration  int
	ReviewDuration  int
	Merged          int
	Closed          int
	Approved        int
	Reviewed        int
	MergedAt        int
	OpenedAt        int
	ClosedAt        int
	LastUpdatedAt   int
	StartedCodingAt int
	StartedPickupAt int
	StartedReviewAt int
}

func CheckValue(t int, max int) int {
	if t > max {
		return 0
	}
	return t
}

func (a *App) Database(c echo.Context) error {

	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)
	page := &templates.Page{
		Title:   "Database",
		Boosted: h.HxBoosted,
	}

	err := godotenv.Load()

	if err != nil {
		return err
	}

	db, err := sql.Open("libsql", os.Getenv("DATABASE_URL"))

	if err != nil {
		return err
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	rows, err := db.Query("SELECT merge_request_metrics.id, merge_request_metrics.users_junk, merge_request_metrics.mr_size, merge_request_metrics.coding_duration, merge_request_metrics.review_start_delay, merge_request_metrics.review_duration, merge_request_metrics.merged, merge_request_metrics.closed, merge_request_metrics.approved, merge_request_metrics.reviewed, merge_request_fact_dates_junk.merged_at, merge_request_fact_dates_junk.opened_at, merge_request_fact_dates_junk.closed_at, merge_request_fact_dates_junk.last_updated_at, merge_request_fact_dates_junk.started_coding_at, merge_request_fact_dates_junk.started_pickup_at, merge_request_fact_dates_junk.started_review_at FROM merge_request_metrics LEFT JOIN merge_request_fact_dates_junk ON merge_request_metrics.dates_junk = merge_request_fact_dates_junk.id;")

	if err != nil {
		return err
	}

	defer rows.Close()

	var joined []JoinedIndexes
	var searchedDates []int
	dateMap := make(map[int]templates.Date)

	for rows.Next() {
		var j JoinedIndexes

		if err := rows.Scan(
			&j.Id,
			&j.UserJunk,
			&j.MrSize,
			&j.CodingDuration,
			&j.PickupDuration,
			&j.ReviewDuration,
			&j.Merged,
			&j.Closed,
			&j.Approved,
			&j.Reviewed,
			&j.MergedAt,
			&j.OpenedAt,
			&j.ClosedAt,
			&j.LastUpdatedAt,
			&j.StartedCodingAt,
			&j.StartedPickupAt,
			&j.StartedReviewAt,
		); err != nil {
			return err
		}
		joined = append(joined, j)
		searchedDates = append(searchedDates, j.MergedAt, j.ClosedAt, j.OpenedAt, j.LastUpdatedAt, j.StartedCodingAt, j.StartedPickupAt, j.StartedReviewAt)
	}

	var metrics []templates.MergeRequestMetrics

	datesQueryString := "SELECT id, day, week, month, year FROM dates WHERE id IN (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(searchedDates)), ", "), "[]") + ");"
	neededDates, er := db.Query(datesQueryString)
	if er != nil {
		return er
	}
	defer neededDates.Close()

	for neededDates.Next() {
		var m templates.Date
		if err := neededDates.Scan(
			&m.Id,
			&m.Day,
			&m.Week,
			&m.Month,
			&m.Year,
		); err != nil {
			return err
		}
		dateMap[m.Id] = templates.Date{
			Day:   CheckValue(m.Day, 31),
			Month: CheckValue(m.Month, 12),
			Year:  CheckValue(m.Year, 2100),
			Week:  CheckValue(m.Week, 52),
		}
	}

	for _, data := range joined {

		metrics = append(metrics, templates.MergeRequestMetrics{
			Id:              data.Id,
			UserJunk:        data.UserJunk,
			MrSize:          data.MrSize,
			CodingDuration:  data.CodingDuration,
			PickupDuration:  data.PickupDuration,
			ReviewDuration:  data.ReviewDuration,
			Merged:          data.Merged,
			Closed:          data.Closed,
			Approved:        data.Approved,
			Reviewed:        data.Reviewed,
			MergedAt:        dateMap[data.MergedAt],
			OpenedAt:        dateMap[data.OpenedAt],
			ClosedAt:        dateMap[data.ClosedAt],
			LastUpdatedAt:   dateMap[data.LastUpdatedAt],
			StartedCodingAt: dateMap[data.StartedCodingAt],
			StartedPickupAt: dateMap[data.StartedPickupAt],
			StartedReviewAt: dateMap[data.StartedReviewAt],
		})
	}

	components := templates.Database(page, page.Title, metrics)
	return components.Render(context.Background(), c.Response().Writer)
}

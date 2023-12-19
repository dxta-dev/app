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
	Author          int
	MergedBy        int
	Approver1       int
	Approver2       int
	Approver3       int
	Approver4       int
	Approver5       int
	Approver6       int
	Approver7       int
	Approver8       int
	Approver9       int
	Approver10      int
	Committer1      int
	Committer2      int
	Committer3      int
	Committer4      int
	Committer5      int
	Committer6      int
	Committer7      int
	Committer8      int
	Committer9      int
	Committer10     int
	Reviewer1       int
	Reviewer2       int
	Reviewer3       int
	Reviewer4       int
	Reviewer5       int
	Reviewer6       int
	Reviewer7       int
	Reviewer8       int
	Reviewer9       int
	Reviewer10      int
}

func CheckValue(t int, max int) int {
	if t > max {
		return 0
	}
	return t
}

func CheckName(n string) string {
	if n == "" {
		return ""
	}
	return n
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

	rows, err := db.Query("SELECT merge_request_metrics.id, merge_request_metrics.merged, merge_request_metrics.closed, merge_request_metrics.approved, merge_request_metrics.reviewed, merge_request_fact_dates_junk.merged_at, merge_request_fact_dates_junk.opened_at, merge_request_fact_dates_junk.closed_at, merge_request_fact_dates_junk.last_updated_at, merge_request_fact_dates_junk.started_coding_at, merge_request_fact_dates_junk.started_pickup_at, merge_request_fact_dates_junk.started_review_at, merge_request_fact_users_junk.author, merge_request_fact_users_junk.merged_by, merge_request_fact_users_junk.approver1, merge_request_fact_users_junk.approver2, merge_request_fact_users_junk.approver3, merge_request_fact_users_junk.approver4, merge_request_fact_users_junk.approver5, merge_request_fact_users_junk.approver6, merge_request_fact_users_junk.approver7, merge_request_fact_users_junk.approver8, merge_request_fact_users_junk.approver9, merge_request_fact_users_junk.approver10, merge_request_fact_users_junk.committer1, merge_request_fact_users_junk.committer2, merge_request_fact_users_junk.committer3, merge_request_fact_users_junk.committer4, merge_request_fact_users_junk.committer5, merge_request_fact_users_junk.committer6, merge_request_fact_users_junk.committer7, merge_request_fact_users_junk.committer8, merge_request_fact_users_junk.committer9, merge_request_fact_users_junk.committer10, merge_request_fact_users_junk.reviewer1, merge_request_fact_users_junk.reviewer2, merge_request_fact_users_junk.reviewer3, merge_request_fact_users_junk.reviewer4, merge_request_fact_users_junk.reviewer5, merge_request_fact_users_junk.reviewer6, merge_request_fact_users_junk.reviewer7, merge_request_fact_users_junk.reviewer8, merge_request_fact_users_junk.reviewer9, merge_request_fact_users_junk.reviewer10 FROM merge_request_metrics LEFT JOIN merge_request_fact_dates_junk ON merge_request_metrics.dates_junk = merge_request_fact_dates_junk.id LEFT JOIN merge_request_fact_users_junk ON merge_request_metrics.users_junk = merge_request_fact_users_junk.id;")

	if err != nil {
		return err
	}

	defer rows.Close()

	var joined []JoinedIndexes
	var searchedDates []int
	var searchedUserJunks []int
	dateMap := make(map[int]templates.Date)
	userMap := make(map[int]templates.User)

	for rows.Next() {
		var j JoinedIndexes

		if err := rows.Scan(
			&j.Id,
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
			&j.Author,
			&j.MergedBy,
			&j.Approver1,
			&j.Approver2,
			&j.Approver3,
			&j.Approver4,
			&j.Approver5,
			&j.Approver6,
			&j.Approver7,
			&j.Approver8,
			&j.Approver9,
			&j.Approver10,
			&j.Committer1,
			&j.Committer2,
			&j.Committer3,
			&j.Committer4,
			&j.Committer5,
			&j.Committer6,
			&j.Committer7,
			&j.Committer8,
			&j.Committer9,
			&j.Committer10,
			&j.Reviewer1,
			&j.Reviewer2,
			&j.Reviewer3,
			&j.Reviewer4,
			&j.Reviewer5,
			&j.Reviewer6,
			&j.Reviewer7,
			&j.Reviewer8,
			&j.Reviewer9,
			&j.Reviewer10,
		); err != nil {
			return err
		}
		joined = append(joined, j)
		searchedDates = append(searchedDates, j.MergedAt, j.ClosedAt, j.OpenedAt, j.LastUpdatedAt, j.StartedCodingAt, j.StartedPickupAt, j.StartedReviewAt)
		searchedUserJunks = append(searchedUserJunks, j.Author, j.MergedBy, j.Approver1, j.Approver2, j.Approver3, j.Approver4, j.Approver5, j.Approver6, j.Approver7, j.Approver8, j.Approver9, j.Approver10, j.Committer1, j.Committer2, j.Committer3, j.Committer4, j.Committer5, j.Committer6, j.Committer7, j.Committer8, j.Committer9, j.Committer10, j.Reviewer1, j.Reviewer2, j.Reviewer3, j.Reviewer4, j.Reviewer5, j.Reviewer6, j.Reviewer7, j.Reviewer8, j.Reviewer9, j.Reviewer10)
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

	usersQueryString := "SELECT id, name FROM forge_users WHERE id IN (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(searchedUserJunks)), ", "), "[]") + ");"
	neededUsers, er := db.Query(usersQueryString)
	if er != nil {
		return er
	}
	defer neededUsers.Close()

	for neededUsers.Next() {
		var m templates.User
		if err := neededUsers.Scan(
			&m.Id,
			&m.Name,
		); err != nil {
			return err
		}
		userMap[m.Id] = templates.User{
			Id:   m.Id,
			Name: CheckName(m.Name),
		}
	}

	for _, data := range joined {

		metrics = append(metrics, templates.MergeRequestMetrics{
			Id:              data.Id,
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
			Author:          userMap[data.Author],
			MergedBy:        userMap[data.MergedBy],
			Approver1:       userMap[data.Approver1],
			Approver2:       userMap[data.Approver2],
			Approver3:       userMap[data.Approver3],
			Approver4:       userMap[data.Approver4],
			Approver5:       userMap[data.Approver5],
			Approver6:       userMap[data.Approver6],
			Approver7:       userMap[data.Approver7],
			Approver8:       userMap[data.Approver8],
			Approver9:       userMap[data.Approver9],
			Approver10:      userMap[data.Approver10],
			Committer1:      userMap[data.Committer1],
			Committer2:      userMap[data.Committer2],
			Committer3:      userMap[data.Committer3],
			Committer4:      userMap[data.Committer4],
			Committer5:      userMap[data.Committer5],
			Committer6:      userMap[data.Committer6],
			Committer7:      userMap[data.Committer7],
			Committer8:      userMap[data.Committer8],
			Committer9:      userMap[data.Committer9],
			Committer10:     userMap[data.Committer10],
			Reviewer1:       userMap[data.Reviewer1],
			Reviewer2:       userMap[data.Reviewer2],
			Reviewer3:       userMap[data.Reviewer3],
			Reviewer4:       userMap[data.Reviewer4],
			Reviewer5:       userMap[data.Reviewer5],
			Reviewer6:       userMap[data.Reviewer6],
			Reviewer7:       userMap[data.Reviewer7],
			Reviewer8:       userMap[data.Reviewer8],
			Reviewer9:       userMap[data.Reviewer9],
			Reviewer10:      userMap[data.Reviewer10],
		})
	}

	components := templates.Database(page, page.Title, metrics)
	return components.Render(context.Background(), c.Response().Writer)
}

package handlers

import (
	"context"
	"database/sql"
	"dxta-dev/app/internals/templates"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type Date struct {
	Day   int
	Month int
	Year  int
}

type User struct {
	Name string
}

type JoinedData struct {
	Id                int
	MergedDate        Date
	OpenedDate        Date
	ClosedDate        Date
	LastUpdatedDate   Date
	StartedCodingDate Date
	StartedPickupDate Date
	StartedReviewDate Date
	Author            User
	MergedBy          User
	Approver1         User
	Approver2         User
	Approver3         User
	Approver4         User
	Approver5         User
	Approver6         User
	Approver7         User
	Approver8         User
	Approver9         User
	Approver10        User
	Committer1        User
	Committer2        User
	Committer3        User
	Committer4        User
	Committer5        User
	Committer6        User
	Committer7        User
	Committer8        User
	Committer9        User
	Committer10       User
	Reviewer1         User
	Reviewer2         User
	Reviewer3         User
	Reviewer4         User
	Reviewer5         User
	Reviewer6         User
	Reviewer7         User
	Reviewer8         User
	Reviewer9         User
	Reviewer10        User
}

func (a *App) Database(c echo.Context) error {

	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	w := r.URL.Query().Get("week")

	fmt.Println("Week: ", w)

	fmt.Println(r.Context().Value(htmx.ContextRequestHeader))

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

	var attributes []interface{}

	// Just for testing
	// Using date 27.07.2023. week 23
	// start of week 24.07.2023.
	// end of week 30.07.2023.

	selectedDate := time.Date(2023, time.Month(7), 27, 0, 0, 0, 0, time.UTC)
	year, week := selectedDate.ISOWeek()

	attributes = append(attributes, week, year)

	selectQuery := fmt.Sprintf(`
		SELECT
		merge_request_metrics.id,
		merged_date.day, merged_date.month, merged_date.year,
		opened_date.day, opened_date.month, opened_date.year,
		closed_date.day, closed_date.month, closed_date.year,
		last_updated_date.day, last_updated_date.month, last_updated_date.year,
		started_coding_date.day, started_coding_date.month, started_coding_date.year,
		started_pickup_date.day, started_pickup_date.month, started_pickup_date.year,
		started_review_date.day, started_review_date.month, started_review_date.year,
		author.name, merged_by.name,
		approver1.name, approver2.name, approver3.name, approver4.name, approver5.name, approver6.name, approver7.name, approver8.name, approver9.name, approver10.name,
		committer1.name, committer2.name, committer3.name, committer4.name, committer5.name, committer6.name, committer7.name, committer8.name, committer9.name, committer10.name,
		reviewer1.name, reviewer2.name, reviewer3.name, reviewer4.name, reviewer5.name, reviewer6.name, reviewer7.name, reviewer8.name, reviewer9.name, reviewer10.name
		FROM merge_request_metrics
		JOIN merge_request_fact_dates_junk
		ON merge_request_metrics.dates_junk = merge_request_fact_dates_junk.id
		JOIN dates AS merged_date
		ON merged_date.id = merge_request_fact_dates_junk.merged_at
		JOIN dates AS opened_date
		ON opened_date.id = merge_request_fact_dates_junk.opened_at
		JOIN dates AS closed_date
		ON closed_date.id = merge_request_fact_dates_junk.closed_at
		JOIN dates AS last_updated_date
		ON last_updated_date.id = merge_request_fact_dates_junk.last_updated_at
		JOIN dates AS started_coding_date
		ON started_coding_date.id = merge_request_fact_dates_junk.started_coding_at
		JOIN dates AS started_pickup_date
		ON started_pickup_date.id = merge_request_fact_dates_junk.started_pickup_at
		JOIN dates AS started_review_date
		ON started_review_date.id = merge_request_fact_dates_junk.started_review_at
		JOIN merge_request_fact_users_junk
		ON merge_request_metrics.users_junk = merge_request_fact_users_junk.id
		JOIN forge_users AS author
		ON author.id = merge_request_fact_users_junk.author
		JOIN forge_users AS merged_by
		ON merged_by.id = merge_request_fact_users_junk.merged_by
		JOIN forge_users AS approver1
		ON approver1.id = merge_request_fact_users_junk.approver1
		JOIN forge_users AS approver2
		ON approver2.id = merge_request_fact_users_junk.approver2
		JOIN forge_users AS approver3
		ON approver3.id = merge_request_fact_users_junk.approver3
		JOIN forge_users AS approver4
		ON approver4.id = merge_request_fact_users_junk.approver4
		JOIN forge_users AS approver5
		ON approver5.id = merge_request_fact_users_junk.approver5
		JOIN forge_users AS approver6
		ON approver6.id = merge_request_fact_users_junk.approver6
		JOIN forge_users AS approver7
		ON approver7.id = merge_request_fact_users_junk.approver7
		JOIN forge_users AS approver8
		ON approver8.id = merge_request_fact_users_junk.approver8
		JOIN forge_users AS approver9
		ON approver9.id = merge_request_fact_users_junk.approver9
		JOIN forge_users AS approver10
		ON approver10.id = merge_request_fact_users_junk.approver10
		JOIN forge_users AS committer1
		ON committer1.id = merge_request_fact_users_junk.committer1
		JOIN forge_users AS committer2
		ON committer2.id = merge_request_fact_users_junk.committer2
		JOIN forge_users AS committer3
		ON committer3.id = merge_request_fact_users_junk.committer3
		JOIN forge_users AS committer4
		ON committer4.id = merge_request_fact_users_junk.committer4
		JOIN forge_users AS committer5
		ON committer5.id = merge_request_fact_users_junk.committer5
		JOIN forge_users AS committer6
		ON committer6.id = merge_request_fact_users_junk.committer6
		JOIN forge_users AS committer7
		ON committer7.id = merge_request_fact_users_junk.committer7
		JOIN forge_users AS committer8
		ON committer8.id = merge_request_fact_users_junk.committer8
		JOIN forge_users AS committer9
		ON committer9.id = merge_request_fact_users_junk.committer9
		JOIN forge_users AS committer10
		ON committer10.id = merge_request_fact_users_junk.committer10
		JOIN forge_users AS reviewer1
		ON reviewer1.id = merge_request_fact_users_junk.reviewer1
		JOIN forge_users AS reviewer2
		ON reviewer2.id = merge_request_fact_users_junk.reviewer2
		JOIN forge_users AS reviewer3
		ON reviewer3.id = merge_request_fact_users_junk.reviewer3
		JOIN forge_users AS reviewer4
		ON reviewer4.id = merge_request_fact_users_junk.reviewer4
		JOIN forge_users AS reviewer5
		ON reviewer5.id = merge_request_fact_users_junk.reviewer5
		JOIN forge_users AS reviewer6
		ON reviewer6.id = merge_request_fact_users_junk.reviewer6
		JOIN forge_users AS reviewer7
		ON reviewer7.id = merge_request_fact_users_junk.reviewer7
		JOIN forge_users AS reviewer8
		ON reviewer8.id = merge_request_fact_users_junk.reviewer8
		JOIN forge_users AS reviewer9
		ON reviewer9.id = merge_request_fact_users_junk.reviewer9
		JOIN forge_users AS reviewer10
		ON reviewer10.id = merge_request_fact_users_junk.reviewer10
		WHERE last_updated_date.week = ?
		AND last_updated_date.year = ?
		;
	`)

	mrStmt, err := db.Prepare(selectQuery)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer mrStmt.Close()

	rows, err := mrStmt.Query(attributes...)

	if err != nil {
		return err
	}

	defer rows.Close()

	var metrics []templates.MergeRequestMetrics
	for rows.Next() {
		var approvers []string
		var committers []string
		var reviewers []string

		data := JoinedData{}
		d := reflect.ValueOf(&data).Elem()
		numCols := d.NumField()
		// Need to find a better wat to get number of columns
		columns := make([]interface{}, 54)

		index := 0
		for i := 0; i < numCols; i++ {
			field := d.Field(i)
			if field.Kind() == reflect.Struct {
				innerData := field
				for j := 0; j < innerData.NumField(); j++ {
					innerField := innerData.Field(j)
					columns[index] = innerField.Addr().Interface()
					index = index + 1
				}
			} else {
				columns[index] = field.Addr().Interface()
				index = index + 1
			}
		}

		if err := rows.Scan(
			columns...,
		); err != nil {
			return err
		}

		approvers = append(approvers, data.Approver1.Name, data.Approver2.Name, data.Approver3.Name, data.Approver4.Name, data.Approver5.Name, data.Approver6.Name, data.Approver7.Name, data.Approver8.Name, data.Approver9.Name, data.Approver10.Name)
		committers = append(committers, data.Committer1.Name, data.Committer2.Name, data.Committer3.Name, data.Committer4.Name, data.Committer5.Name, data.Committer6.Name, data.Committer7.Name, data.Committer8.Name, data.Committer9.Name, data.Committer10.Name)
		reviewers = append(reviewers, data.Reviewer1.Name, data.Reviewer2.Name, data.Reviewer3.Name, data.Reviewer4.Name, data.Reviewer5.Name, data.Reviewer6.Name, data.Reviewer7.Name, data.Reviewer8.Name, data.Reviewer9.Name, data.Reviewer10.Name)

		metrics = append(metrics, templates.MergeRequestMetrics{
			Id:              data.Id,
			MergedAt:        time.Date(data.MergedDate.Year, time.Month(data.MergedDate.Month), data.MergedDate.Day, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			OpenedAt:        time.Date(data.OpenedDate.Year, time.Month(data.OpenedDate.Month), data.OpenedDate.Day, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			ClosedAt:        time.Date(data.ClosedDate.Year, time.Month(data.OpenedDate.Month), data.ClosedDate.Day, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			LastUpdatedAt:   time.Date(data.LastUpdatedDate.Year, time.Month(data.LastUpdatedDate.Month), data.LastUpdatedDate.Day, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			StartedCodingAt: time.Date(data.StartedCodingDate.Year, time.Month(data.StartedCodingDate.Month), data.StartedCodingDate.Day, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			StartedPickupAt: time.Date(data.StartedPickupDate.Year, time.Month(data.StartedPickupDate.Month), data.StartedPickupDate.Day, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			StartedReviewAt: time.Date(data.StartedReviewDate.Year, time.Month(data.StartedReviewDate.Month), data.StartedReviewDate.Day, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			Author:          data.Author.Name,
			MergedBy:        data.MergedBy.Name,
			Approvers:       approvers,
			Committers:      committers,
			Reviewers:       reviewers,
		})
	}

	if(w != "") {
		res := c.Response()

		res.Header().Set("HX-Push-Url", "/database/" + w)
	}

	components := templates.Database(page, page.Title, metrics)
	return components.Render(context.Background(), c.Response().Writer)
}

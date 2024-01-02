package handlers

import (
	"context"
	"database/sql"
	"dxta-dev/app/internals/templates"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type JoinedIndexes struct {
	Id                 int
	MergedDay          int
	MergedMonth        int
	MergedYear         int
	OpenedDay          int
	OpenedMonth        int
	OpenedYear         int
	ClosedDay          int
	ClosedMonth        int
	ClosedYear         int
	LastUpdatedDay     int
	LastUpdatedMonth   int
	LastUpdatedYear    int
	StartedCodingDay   int
	StartedCodingMonth int
	StartedCodingYear  int
	StartedPickupDay   int
	StartedPickupMonth int
	StartedPickupYear  int
	StartedReviewDay   int
	StartedReviewMonth int
	StartedReviewYear  int
	AuthorName         string
	MergedByName       string
	Approver1Name      string
	Approver2Name      string
	Approver3Name      string
	Approver4Name      string
	Approver5Name      string
	Approver6Name      string
	Approver7Name      string
	Approver8Name      string
	Approver9Name      string
	Approver10Name     string
	Committer1Name     string
	Committer2Name     string
	Committer3Name     string
	Committer4Name     string
	Committer5Name     string
	Committer6Name     string
	Committer7Name     string
	Committer8Name     string
	Committer9Name     string
	Committer10Name    string
	Reviewer1Name      string
	Reviewer2Name      string
	Reviewer3Name      string
	Reviewer4Name      string
	Reviewer5Name      string
	Reviewer6Name      string
	Reviewer7Name      string
	Reviewer8Name      string
	Reviewer9Name      string
	Reviewer10Name     string
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

	// Just for testing
	repositoryId := 2
	mrIds := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31}
	var mrIdsInterface []interface{}

	placeholderMrIdsSplice := make([]string, len(mrIds))
	for i := range placeholderMrIdsSplice {
		placeholderMrIdsSplice[i] = "?"
		mrIdsInterface = append(mrIdsInterface, mrIds[i])
	}
	placeholderMrIds := strings.Join(placeholderMrIdsSplice, ", ")

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
		WHERE merge_request_metrics.merge_request IN (%s) 
		AND merge_request_metrics.repository = %d
		;
	`, placeholderMrIds, repositoryId)

	mrStmt, err := db.Prepare(selectQuery)
	if err != nil {
		fmt.Println("US Error:", err)
		return err
	}
	defer mrStmt.Close()

	rows, err := mrStmt.Query(mrIdsInterface...)

	if err != nil {
		return err
	}

	defer rows.Close()

	var metrics []templates.MergeRequestMetrics
	for rows.Next() {
		var approvers []string
		var committers []string
		var reviewers []string
		var j JoinedIndexes

		if err := rows.Scan(
			&j.Id,
			&j.MergedDay,
			&j.MergedMonth,
			&j.MergedYear,
			&j.OpenedDay,
			&j.OpenedMonth,
			&j.OpenedYear,
			&j.ClosedDay,
			&j.ClosedMonth,
			&j.ClosedYear,
			&j.LastUpdatedDay,
			&j.LastUpdatedMonth,
			&j.LastUpdatedYear,
			&j.StartedCodingDay,
			&j.StartedCodingMonth,
			&j.StartedCodingYear,
			&j.StartedPickupDay,
			&j.StartedPickupMonth,
			&j.StartedPickupYear,
			&j.StartedReviewDay,
			&j.StartedReviewMonth,
			&j.StartedReviewYear,
			&j.AuthorName,
			&j.MergedByName,
			&j.Approver1Name,
			&j.Approver2Name,
			&j.Approver3Name,
			&j.Approver4Name,
			&j.Approver5Name,
			&j.Approver6Name,
			&j.Approver7Name,
			&j.Approver8Name,
			&j.Approver9Name,
			&j.Approver10Name,
			&j.Committer1Name,
			&j.Committer2Name,
			&j.Committer3Name,
			&j.Committer4Name,
			&j.Committer5Name,
			&j.Committer6Name,
			&j.Committer7Name,
			&j.Committer8Name,
			&j.Committer9Name,
			&j.Committer10Name,
			&j.Reviewer1Name,
			&j.Reviewer2Name,
			&j.Reviewer3Name,
			&j.Reviewer4Name,
			&j.Reviewer5Name,
			&j.Reviewer6Name,
			&j.Reviewer7Name,
			&j.Reviewer8Name,
			&j.Reviewer9Name,
			&j.Reviewer10Name,
		); err != nil {
			return err
		}

		approvers = append(approvers, j.Approver1Name, j.Approver2Name, j.Approver3Name, j.Approver4Name, j.Approver5Name, j.Approver6Name, j.Approver7Name, j.Approver8Name, j.Approver9Name, j.Approver10Name)
		committers = append(committers, j.Committer1Name, j.Committer2Name, j.Committer3Name, j.Committer4Name, j.Committer5Name, j.Committer6Name, j.Committer7Name, j.Committer8Name, j.Committer9Name, j.Committer10Name)
		reviewers = append(reviewers, j.Reviewer1Name, j.Reviewer2Name, j.Reviewer3Name, j.Reviewer4Name, j.Reviewer5Name, j.Reviewer6Name, j.Reviewer7Name, j.Reviewer8Name, j.Reviewer9Name, j.Reviewer10Name)

		metrics = append(metrics, templates.MergeRequestMetrics{
			Id:              j.Id,
			MergedAt:        time.Date(j.MergedYear, time.Month(j.MergedMonth), j.MergedDay, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			OpenedAt:        time.Date(j.OpenedYear, time.Month(j.OpenedMonth), j.OpenedDay, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			ClosedAt:        time.Date(j.ClosedYear, time.Month(j.MergedMonth), j.ClosedDay, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			LastUpdatedAt:   time.Date(j.LastUpdatedYear, time.Month(j.LastUpdatedMonth), j.LastUpdatedDay, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			StartedCodingAt: time.Date(j.StartedCodingYear, time.Month(j.StartedCodingMonth), j.StartedCodingDay, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			StartedPickupAt: time.Date(j.StartedPickupYear, time.Month(j.StartedPickupMonth), j.StartedPickupDay, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			StartedReviewAt: time.Date(j.StartedReviewYear, time.Month(j.StartedReviewMonth), j.StartedReviewDay, 0, 0, 0, 0, time.UTC).Format("Mon, 02-01-2006"),
			Author:          j.AuthorName,
			MergedBy:        j.MergedByName,
			Approvers:       approvers,
			Committers:      reviewers,
			Reviewers:       reviewers,
		})
	}

	components := templates.Database(page, page.Title, metrics)
	return components.Render(context.Background(), c.Response().Writer)
}

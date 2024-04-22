package data

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dxta-dev/app/internal/util"
)

type MergeRequestListItemData struct {
	Id            int64
	Title         string
	WebUrl        string
	CanonId       int64
	CodeAdditions int64
	CodeDeletions int64
	LastUpdatedAt time.Time
}

func (s *Store) GetMergeRequestsInProgress(date time.Time, teamMembers []int64) ([]MergeRequestListItemData, error) {
	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, err
	}

	week := util.GetFormattedWeek(date)

	query := fmt.Sprintf(`
	SELECT DISTINCT
		mr.id,
		mr.title,
		mr.web_url,
		mr.canon_id,
		metrics.code_addition,
		metrics.code_deletion,
		last_updated_at.day,
		last_updated_at.month,
		last_updated_at.year
	FROM transform_merge_request_events AS events
	JOIN transform_dates AS occured_on ON occured_on.id = events.occured_on
	JOIN transform_forge_users AS user ON user.id = events.actor
	JOIN transform_merge_requests AS mr ON mr.id = events.merge_request
	JOIN transform_merge_request_metrics AS metrics ON metrics.merge_request = mr.id
	JOIN transform_merge_request_fact_users_junk AS u ON u.id = metrics.users_junk
	JOIN transform_merge_request_fact_dates_junk AS dj ON dj.id = metrics.dates_junk
	JOIN transform_dates AS last_updated_at ON last_updated_at.id = dj.last_updated_at
	JOIN transform_forge_users AS author ON author.id = u.author
	WHERE
		occured_on.week = ?
		AND events.merge_request_event_type = 9
  	AND metrics.reviewed = 0
  	AND metrics.approved = 0
  	AND metrics.merged = 0
  	AND metrics.closed = 0
		AND author.bot = 0
		AND user.bot = 0
	%s
	LIMIT 5
	`, usersInTeamConditionQuery)

	defer db.Close()

	queryParams := make([]interface{}, len(teamMembers)+1)
	queryParams[0] = week
	for i, v := range teamMembers {
		queryParams[i+1] = v
	}

	rows, err := db.Query(query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var mergeRequests []MergeRequestListItemData

	for rows.Next() {
		var item MergeRequestListItemData
		var lastUpdatedAtDate int64
		var lastUpdatedAtMonth int64
		var lastUpdatedAtYear int64

		if err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.WebUrl,
			&item.CanonId,
			&item.CodeAdditions,
			&item.CodeDeletions,
			&lastUpdatedAtDate,
			&lastUpdatedAtMonth,
			&lastUpdatedAtYear); err != nil {
			return nil, err
		}

		item.LastUpdatedAt = time.Date(int(lastUpdatedAtYear), time.Month(lastUpdatedAtMonth), int(lastUpdatedAtDate), 0, 0, 0, 0, time.UTC)

		mergeRequests = append(mergeRequests, item)
	}

	return mergeRequests, nil
}

func (s *Store) GetMergeRequests(date time.Time, teamMembers []int64) (interface{}, error) {
	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	week := util.GetFormattedWeek(date)

	query := fmt.Sprintf(`
	SELECT
		mr.id,
		mr.title,
		mr.web_url,
		mr.canon_id,
		metrics.code_addition,
		metrics.code_deletion,
		metrics.merged,
		metrics.closed,
		metrics.approved,
		metrics.reviewed
	FROM transform_merge_request_events AS ev
	JOIN transform_dates AS date ON date.id = ev.occured_on
	JOIN transform_forge_users AS user ON user.id = ev.actor
	JOIN transform_merge_requests AS mr ON mr.id = ev.merge_request
	JOIN transform_merge_request_metrics AS metrics ON metrics.merge_request = mr.id
	JOIN transform_merge_request_fact_users_junk AS u ON u.id = metrics.users_junk
	JOIN transform_forge_users AS author ON author.id = u.author
	WHERE date.week = '2024-W13'
	AND author.bot = 0
	AND user.bot = 0
	%s
	GROUP BY mr.id;
	`, usersInTeamConditionQuery)

	queryParams := make([]interface{}, len(teamMembers)+1)
	queryParams[0] = week
	for i, v := range teamMembers {
		queryParams[i+1] = v
	}

	_, err = db.Query(query, queryParams...)
	return nil, nil
}

package data

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dxta-dev/app/internal/util"
)

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

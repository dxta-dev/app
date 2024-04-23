package data

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dxta-dev/app/internal/util"
)

type UserAvatarUrl struct {
	UserId int64
	Url    string
}

type MergeRequestListItemData struct {
	Id             int64
	Title          string
	WebUrl         string
	CanonId        int64
	CodeAdditions  int64
	CodeDeletions  int64
	LastUpdatedAt  time.Time
	UserAvatarUrls []string
}

const iMAX_USER_AVATARS_LEN = 6

func (s *Store) GetMergeRequestsInProgress(date time.Time, teamMembers []int64, nullUserId int64) ([]MergeRequestListItemData, error) {
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
		last_updated_at.year,
		author.id, author.avatar_url,
		merger.id, merger.avatar_url,
		approver1.id,  approver1.avatar_url,
		approver2.id,  approver2.avatar_url,
		approver3.id,  approver3.avatar_url,
		approver4.id,  approver4.avatar_url,
		approver5.id,  approver5.avatar_url,
		approver6.id,  approver6.avatar_url,
		approver7.id,  approver7.avatar_url,
		approver8.id,  approver8.avatar_url,
		approver9.id,  approver9.avatar_url,
		approver10.id, approver10.avatar_url,
		committer1.id,  committer1.avatar_url,
		committer2.id,  committer2.avatar_url,
		committer3.id,  committer3.avatar_url,
		committer4.id,  committer4.avatar_url,
		committer5.id,  committer5.avatar_url,
		committer6.id,  committer6.avatar_url,
		committer7.id,  committer7.avatar_url,
		committer8.id,  committer8.avatar_url,
		committer9.id,  committer9.avatar_url,
		committer10.id, committer10.avatar_url,
		reviewer1.id,  reviewer1.avatar_url,
		reviewer2.id,  reviewer2.avatar_url,
		reviewer3.id,  reviewer3.avatar_url,
		reviewer4.id,  reviewer4.avatar_url,
		reviewer5.id,  reviewer5.avatar_url,
		reviewer6.id,  reviewer6.avatar_url,
		reviewer7.id,  reviewer7.avatar_url,
		reviewer8.id,  reviewer8.avatar_url,
		reviewer9.id,  reviewer9.avatar_url,
		reviewer10.id, reviewer10.avatar_url
	FROM transform_merge_request_events AS events
	JOIN transform_dates AS occured_on ON occured_on.id = events.occured_on
	JOIN transform_forge_users AS user ON user.id = events.actor
	JOIN transform_merge_requests AS mr ON mr.id = events.merge_request
	JOIN transform_merge_request_metrics AS metrics ON metrics.merge_request = mr.id
	JOIN transform_merge_request_fact_users_junk AS u ON u.id = metrics.users_junk
	JOIN transform_merge_request_fact_dates_junk AS dj ON dj.id = metrics.dates_junk
	JOIN transform_dates AS last_updated_at ON last_updated_at.id = dj.last_updated_at
	JOIN transform_forge_users AS author ON author.id = u.author
	JOIN transform_forge_users AS merger ON merger.id = u.merged_by
	JOIN transform_forge_users AS approver1   ON approver1.id  = u.approver1
	JOIN transform_forge_users AS approver2   ON approver2.id  = u.approver2
	JOIN transform_forge_users AS approver3   ON approver3.id  = u.approver3
	JOIN transform_forge_users AS approver4   ON approver4.id  = u.approver4
	JOIN transform_forge_users AS approver5   ON approver5.id  = u.approver5
	JOIN transform_forge_users AS approver6   ON approver6.id  = u.approver6
	JOIN transform_forge_users AS approver7   ON approver7.id  = u.approver7
	JOIN transform_forge_users AS approver8   ON approver8.id  = u.approver8
	JOIN transform_forge_users AS approver9   ON approver9.id  = u.approver9
	JOIN transform_forge_users AS approver10  ON approver10.id = u.approver10
	JOIN transform_forge_users AS committer1   ON committer1.id  = u.committer1
	JOIN transform_forge_users AS committer2   ON committer2.id  = u.committer2
	JOIN transform_forge_users AS committer3   ON committer3.id  = u.committer3
	JOIN transform_forge_users AS committer4   ON committer4.id  = u.committer4
	JOIN transform_forge_users AS committer5   ON committer5.id  = u.committer5
	JOIN transform_forge_users AS committer6   ON committer6.id  = u.committer6
	JOIN transform_forge_users AS committer7   ON committer7.id  = u.committer7
	JOIN transform_forge_users AS committer8   ON committer8.id  = u.committer8
	JOIN transform_forge_users AS committer9   ON committer9.id  = u.committer9
	JOIN transform_forge_users AS committer10 ON committer10.id = u.committer10
	JOIN transform_forge_users AS reviewer1   ON reviewer1.id  = u.reviewer1
	JOIN transform_forge_users AS reviewer2   ON reviewer2.id  = u.reviewer2
	JOIN transform_forge_users AS reviewer3   ON reviewer3.id  = u.reviewer3
	JOIN transform_forge_users AS reviewer4   ON reviewer4.id  = u.reviewer4
	JOIN transform_forge_users AS reviewer5   ON reviewer5.id  = u.reviewer5
	JOIN transform_forge_users AS reviewer6   ON reviewer6.id  = u.reviewer6
	JOIN transform_forge_users AS reviewer7   ON reviewer7.id  = u.reviewer7
	JOIN transform_forge_users AS reviewer8   ON reviewer8.id  = u.reviewer8
	JOIN transform_forge_users AS reviewer9   ON reviewer9.id  = u.reviewer9
	JOIN transform_forge_users AS reviewer10  ON reviewer10.id = u.reviewer10
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
		var userAvatars = make([]UserAvatarUrl, 2+3*10)

		if err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.WebUrl,
			&item.CanonId,
			&item.CodeAdditions,
			&item.CodeDeletions,
			&lastUpdatedAtDate,
			&lastUpdatedAtMonth,
			&lastUpdatedAtYear,
			&userAvatars[0].UserId, &userAvatars[0].Url,
			&userAvatars[1].UserId, &userAvatars[1].Url,
			&userAvatars[2].UserId, &userAvatars[2].Url,
			&userAvatars[3].UserId, &userAvatars[3].Url,
			&userAvatars[4].UserId, &userAvatars[4].Url,
			&userAvatars[5].UserId, &userAvatars[5].Url,
			&userAvatars[6].UserId, &userAvatars[6].Url,
			&userAvatars[7].UserId, &userAvatars[7].Url,
			&userAvatars[8].UserId, &userAvatars[8].Url,
			&userAvatars[9].UserId, &userAvatars[9].Url,
			&userAvatars[10].UserId, &userAvatars[10].Url,
			&userAvatars[11].UserId, &userAvatars[11].Url,
			&userAvatars[12].UserId, &userAvatars[12].Url,
			&userAvatars[13].UserId, &userAvatars[13].Url,
			&userAvatars[14].UserId, &userAvatars[14].Url,
			&userAvatars[15].UserId, &userAvatars[15].Url,
			&userAvatars[16].UserId, &userAvatars[16].Url,
			&userAvatars[17].UserId, &userAvatars[17].Url,
			&userAvatars[18].UserId, &userAvatars[18].Url,
			&userAvatars[19].UserId, &userAvatars[19].Url,
			&userAvatars[20].UserId, &userAvatars[20].Url,
			&userAvatars[21].UserId, &userAvatars[21].Url,
			&userAvatars[22].UserId, &userAvatars[22].Url,
			&userAvatars[23].UserId, &userAvatars[23].Url,
			&userAvatars[24].UserId, &userAvatars[24].Url,
			&userAvatars[25].UserId, &userAvatars[25].Url,
			&userAvatars[26].UserId, &userAvatars[26].Url,
			&userAvatars[27].UserId, &userAvatars[27].Url,
			&userAvatars[28].UserId, &userAvatars[28].Url,
			&userAvatars[29].UserId, &userAvatars[29].Url,
			&userAvatars[30].UserId, &userAvatars[30].Url,
			&userAvatars[31].UserId, &userAvatars[31].Url,
		); err != nil {
			return nil, err
		}

		uniqueUsersMap := make(map[int64]bool)
		for _, userAvatar := range userAvatars {
			if userAvatar.UserId != nullUserId && !uniqueUsersMap[userAvatar.UserId] {
				uniqueUsersMap[userAvatar.UserId] = true
				item.UserAvatarUrls = append(item.UserAvatarUrls, userAvatar.Url)
			}
		}

		item.LastUpdatedAt = time.Date(int(lastUpdatedAtYear), time.Month(lastUpdatedAtMonth), int(lastUpdatedAtDate), 0, 0, 0, 0, time.UTC)

		mergeRequests = append(mergeRequests, item)
	}

	return mergeRequests, nil
}

func (s *Store) GetMergeRequestsClosed(date time.Time, teamMembers []int64, nullUserId int64, andMerged bool) ([]MergeRequestListItemData, error) {
	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	mergedClauseCondition := "0"
	if andMerged {
		mergedClauseCondition = "1"
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
		last_updated_at.year,
		author.id, author.avatar_url,
		merger.id, merger.avatar_url,
		approver1.id,  approver1.avatar_url,
		approver2.id,  approver2.avatar_url,
		approver3.id,  approver3.avatar_url,
		approver4.id,  approver4.avatar_url,
		approver5.id,  approver5.avatar_url,
		approver6.id,  approver6.avatar_url,
		approver7.id,  approver7.avatar_url,
		approver8.id,  approver8.avatar_url,
		approver9.id,  approver9.avatar_url,
		approver10.id, approver10.avatar_url,
		committer1.id,  committer1.avatar_url,
		committer2.id,  committer2.avatar_url,
		committer3.id,  committer3.avatar_url,
		committer4.id,  committer4.avatar_url,
		committer5.id,  committer5.avatar_url,
		committer6.id,  committer6.avatar_url,
		committer7.id,  committer7.avatar_url,
		committer8.id,  committer8.avatar_url,
		committer9.id,  committer9.avatar_url,
		committer10.id, committer10.avatar_url,
		reviewer1.id,  reviewer1.avatar_url,
		reviewer2.id,  reviewer2.avatar_url,
		reviewer3.id,  reviewer3.avatar_url,
		reviewer4.id,  reviewer4.avatar_url,
		reviewer5.id,  reviewer5.avatar_url,
		reviewer6.id,  reviewer6.avatar_url,
		reviewer7.id,  reviewer7.avatar_url,
		reviewer8.id,  reviewer8.avatar_url,
		reviewer9.id,  reviewer9.avatar_url,
		reviewer10.id, reviewer10.avatar_url
	FROM transform_merge_request_events AS events
	JOIN transform_dates AS occured_on ON occured_on.id = events.occured_on
	JOIN transform_forge_users AS user ON user.id = events.actor
	JOIN transform_merge_requests AS mr ON mr.id = events.merge_request
	JOIN transform_merge_request_metrics AS metrics ON metrics.merge_request = mr.id
	JOIN transform_merge_request_fact_users_junk AS u ON u.id = metrics.users_junk
	JOIN transform_merge_request_fact_dates_junk AS dj ON dj.id = metrics.dates_junk
	JOIN transform_dates AS last_updated_at ON last_updated_at.id = dj.last_updated_at
	JOIN transform_forge_users AS author ON author.id = u.author
	JOIN transform_forge_users AS merger ON merger.id = u.merged_by
	JOIN transform_forge_users AS approver1   ON approver1.id  = u.approver1
	JOIN transform_forge_users AS approver2   ON approver2.id  = u.approver2
	JOIN transform_forge_users AS approver3   ON approver3.id  = u.approver3
	JOIN transform_forge_users AS approver4   ON approver4.id  = u.approver4
	JOIN transform_forge_users AS approver5   ON approver5.id  = u.approver5
	JOIN transform_forge_users AS approver6   ON approver6.id  = u.approver6
	JOIN transform_forge_users AS approver7   ON approver7.id  = u.approver7
	JOIN transform_forge_users AS approver8   ON approver8.id  = u.approver8
	JOIN transform_forge_users AS approver9   ON approver9.id  = u.approver9
	JOIN transform_forge_users AS approver10  ON approver10.id = u.approver10
	JOIN transform_forge_users AS committer1   ON committer1.id  = u.committer1
	JOIN transform_forge_users AS committer2   ON committer2.id  = u.committer2
	JOIN transform_forge_users AS committer3   ON committer3.id  = u.committer3
	JOIN transform_forge_users AS committer4   ON committer4.id  = u.committer4
	JOIN transform_forge_users AS committer5   ON committer5.id  = u.committer5
	JOIN transform_forge_users AS committer6   ON committer6.id  = u.committer6
	JOIN transform_forge_users AS committer7   ON committer7.id  = u.committer7
	JOIN transform_forge_users AS committer8   ON committer8.id  = u.committer8
	JOIN transform_forge_users AS committer9   ON committer9.id  = u.committer9
	JOIN transform_forge_users AS committer10 ON committer10.id = u.committer10
	JOIN transform_forge_users AS reviewer1   ON reviewer1.id  = u.reviewer1
	JOIN transform_forge_users AS reviewer2   ON reviewer2.id  = u.reviewer2
	JOIN transform_forge_users AS reviewer3   ON reviewer3.id  = u.reviewer3
	JOIN transform_forge_users AS reviewer4   ON reviewer4.id  = u.reviewer4
	JOIN transform_forge_users AS reviewer5   ON reviewer5.id  = u.reviewer5
	JOIN transform_forge_users AS reviewer6   ON reviewer6.id  = u.reviewer6
	JOIN transform_forge_users AS reviewer7   ON reviewer7.id  = u.reviewer7
	JOIN transform_forge_users AS reviewer8   ON reviewer8.id  = u.reviewer8
	JOIN transform_forge_users AS reviewer9   ON reviewer9.id  = u.reviewer9
	JOIN transform_forge_users AS reviewer10  ON reviewer10.id = u.reviewer10
	WHERE
		occured_on.week = ?
		AND events.merge_request_event_type = 11
  	AND metrics.merged = %s
  	AND metrics.closed = 1
		AND author.bot = 0
		AND user.bot = 0
	%s
	ORDER BY
  	last_updated_at.year DESC,
  	last_updated_at.month DESC,
  	last_updated_at.day DESC
	LIMIT 5
	`, mergedClauseCondition, usersInTeamConditionQuery)

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
		var userAvatars = make([]UserAvatarUrl, 2+3*10)

		if err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.WebUrl,
			&item.CanonId,
			&item.CodeAdditions,
			&item.CodeDeletions,
			&lastUpdatedAtDate,
			&lastUpdatedAtMonth,
			&lastUpdatedAtYear,
			&userAvatars[0].UserId, &userAvatars[0].Url,
			&userAvatars[1].UserId, &userAvatars[1].Url,
			&userAvatars[2].UserId, &userAvatars[2].Url,
			&userAvatars[3].UserId, &userAvatars[3].Url,
			&userAvatars[4].UserId, &userAvatars[4].Url,
			&userAvatars[5].UserId, &userAvatars[5].Url,
			&userAvatars[6].UserId, &userAvatars[6].Url,
			&userAvatars[7].UserId, &userAvatars[7].Url,
			&userAvatars[8].UserId, &userAvatars[8].Url,
			&userAvatars[9].UserId, &userAvatars[9].Url,
			&userAvatars[10].UserId, &userAvatars[10].Url,
			&userAvatars[11].UserId, &userAvatars[11].Url,
			&userAvatars[12].UserId, &userAvatars[12].Url,
			&userAvatars[13].UserId, &userAvatars[13].Url,
			&userAvatars[14].UserId, &userAvatars[14].Url,
			&userAvatars[15].UserId, &userAvatars[15].Url,
			&userAvatars[16].UserId, &userAvatars[16].Url,
			&userAvatars[17].UserId, &userAvatars[17].Url,
			&userAvatars[18].UserId, &userAvatars[18].Url,
			&userAvatars[19].UserId, &userAvatars[19].Url,
			&userAvatars[20].UserId, &userAvatars[20].Url,
			&userAvatars[21].UserId, &userAvatars[21].Url,
			&userAvatars[22].UserId, &userAvatars[22].Url,
			&userAvatars[23].UserId, &userAvatars[23].Url,
			&userAvatars[24].UserId, &userAvatars[24].Url,
			&userAvatars[25].UserId, &userAvatars[25].Url,
			&userAvatars[26].UserId, &userAvatars[26].Url,
			&userAvatars[27].UserId, &userAvatars[27].Url,
			&userAvatars[28].UserId, &userAvatars[28].Url,
			&userAvatars[29].UserId, &userAvatars[29].Url,
			&userAvatars[30].UserId, &userAvatars[30].Url,
			&userAvatars[31].UserId, &userAvatars[31].Url,
		); err != nil {
			return nil, err
		}

		uniqueUsersMap := make(map[int64]bool)
		for _, userAvatar := range userAvatars {
			if userAvatar.UserId != nullUserId && !uniqueUsersMap[userAvatar.UserId] {
				uniqueUsersMap[userAvatar.UserId] = true
				item.UserAvatarUrls = append(item.UserAvatarUrls, userAvatar.Url)
			}
			if len(item.UserAvatarUrls) >= iMAX_USER_AVATARS_LEN {
				break
			}
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

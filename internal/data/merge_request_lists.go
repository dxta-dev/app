package data

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dxta-dev/app/internal/util"
)

const iMAX_USER_AVATARS_LEN = 6

type RepositoryInfo struct {
	Name string
	Org  string
}

type ListUserInfo struct {
	UserId int64
	Url    string
	Name   string
	Bot    bool
}

type MergeRequestListItemData struct {
	Id                      int64
	Repo                    RepositoryInfo
	Title                   string
	WebUrl                  string
	CanonId                 int64
	CodeAdditions           int64
	CodeDeletions           int64
	ReviewDepth             int64
	LastEventAt             time.Time
	LastSwarmEventAt        time.Time
	LastEventTimestamp      int64
	CurrentWeekState        string
	CurrentMinimapIndicator string
	Actors                  []ListUserInfo
}

const mrListDataSelect = `mr.id,
	mr.title,
	mr.web_url,
	mr.canon_id,
	metrics.code_addition,
	metrics.code_deletion,
	metrics.review_depth,
	repos.name,
	repos.namespace_name,
	MAX(events.timestamp) as last_event_ts,
	author.id, author.avatar_url, author.name, author.bot,
	merger.id, merger.avatar_url, merger.name, merger.bot,
	approver1.id,  approver1.avatar_url,  approver1.name,  approver1.bot,
	approver2.id,  approver2.avatar_url,  approver2.name,  approver2.bot,
	approver3.id,  approver3.avatar_url,  approver3.name,  approver3.bot,
	approver4.id,  approver4.avatar_url,  approver4.name,  approver4.bot,
	approver5.id,  approver5.avatar_url,  approver5.name,  approver5.bot,
	approver6.id,  approver6.avatar_url,  approver6.name,  approver6.bot,
	approver7.id,  approver7.avatar_url,  approver7.name,  approver7.bot,
	approver8.id,  approver8.avatar_url,  approver8.name,  approver8.bot,
	approver9.id,  approver9.avatar_url,  approver9.name,  approver9.bot,
	approver10.id, approver10.avatar_url, approver10.name, approver10.bot,
	committer1.id,  committer1.avatar_url,  committer1.name,  committer1.bot,
	committer2.id,  committer2.avatar_url,  committer2.name,  committer2.bot,
	committer3.id,  committer3.avatar_url,  committer3.name,  committer3.bot,
	committer4.id,  committer4.avatar_url,  committer4.name,  committer4.bot,
	committer5.id,  committer5.avatar_url,  committer5.name,  committer5.bot,
	committer6.id,  committer6.avatar_url,  committer6.name,  committer6.bot,
	committer7.id,  committer7.avatar_url,  committer7.name,  committer7.bot,
	committer8.id,  committer8.avatar_url,  committer8.name,  committer8.bot,
	committer9.id,  committer9.avatar_url,  committer9.name,  committer9.bot,
	committer10.id, committer10.avatar_url, committer10.name, committer10.bot,
	reviewer1.id,  reviewer1.avatar_url,  reviewer1.name,  reviewer1.bot,
	reviewer2.id,  reviewer2.avatar_url,  reviewer2.name,  reviewer2.bot,
	reviewer3.id,  reviewer3.avatar_url,  reviewer3.name,  reviewer3.bot,
	reviewer4.id,  reviewer4.avatar_url,  reviewer4.name,  reviewer4.bot,
	reviewer5.id,  reviewer5.avatar_url,  reviewer5.name,  reviewer5.bot,
	reviewer6.id,  reviewer6.avatar_url,  reviewer6.name,  reviewer6.bot,
	reviewer7.id,  reviewer7.avatar_url,  reviewer7.name,  reviewer7.bot,
	reviewer8.id,  reviewer8.avatar_url,  reviewer8.name,  reviewer8.bot,
	reviewer9.id,  reviewer9.avatar_url,  reviewer9.name,  reviewer9.bot,
	reviewer10.id, reviewer10.avatar_url, reviewer10.name, reviewer10.bot`

const mrListTables = `transform_merge_request_events AS events
	JOIN transform_dates AS occured_on ON occured_on.id = events.occured_on
	JOIN transform_forge_users AS user ON user.id = events.actor
	JOIN transform_merge_requests AS mr ON mr.id = events.merge_request
	JOIN transform_merge_request_metrics AS metrics ON metrics.merge_request = mr.id
	JOIN transform_repositories as repos ON repos.id = metrics.repository
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
	JOIN transform_forge_users AS reviewer10  ON reviewer10.id = u.reviewer10`

func compareWeeks(week string, lastEventDate time.Time) string {

	firstWeekDate, lastWeekDate, err := util.ParseYearWeek(week)

	if err != nil {
		return ""
	}

	if !lastEventDate.Before(firstWeekDate) && !lastEventDate.After(lastWeekDate) {
		return "curr"
	} else if firstWeekDate.After(lastEventDate) {
		return "prev"
	} else {
		return "next"
	}
}

func scanMergeRequestListItemRow(item *MergeRequestListItemData, userAvatars []ListUserInfo, rows *sql.Rows, week string) error {
	var lastEventMilli int64

	err := rows.Scan(
		&item.Id,
		&item.Title,
		&item.WebUrl,
		&item.CanonId,
		&item.CodeAdditions,
		&item.CodeDeletions,
		&item.ReviewDepth,
		&item.Repo.Name,
		&item.Repo.Org,
		&lastEventMilli,
		&userAvatars[0].UserId, &userAvatars[0].Url, &userAvatars[0].Name, &userAvatars[0].Bot,
		&userAvatars[1].UserId, &userAvatars[1].Url, &userAvatars[1].Name, &userAvatars[1].Bot,
		&userAvatars[2].UserId, &userAvatars[2].Url, &userAvatars[2].Name, &userAvatars[2].Bot,
		&userAvatars[3].UserId, &userAvatars[3].Url, &userAvatars[3].Name, &userAvatars[3].Bot,
		&userAvatars[4].UserId, &userAvatars[4].Url, &userAvatars[4].Name, &userAvatars[4].Bot,
		&userAvatars[5].UserId, &userAvatars[5].Url, &userAvatars[5].Name, &userAvatars[5].Bot,
		&userAvatars[6].UserId, &userAvatars[6].Url, &userAvatars[6].Name, &userAvatars[6].Bot,
		&userAvatars[7].UserId, &userAvatars[7].Url, &userAvatars[7].Name, &userAvatars[7].Bot,
		&userAvatars[8].UserId, &userAvatars[8].Url, &userAvatars[8].Name, &userAvatars[8].Bot,
		&userAvatars[9].UserId, &userAvatars[9].Url, &userAvatars[9].Name, &userAvatars[9].Bot,
		&userAvatars[10].UserId, &userAvatars[10].Url, &userAvatars[10].Name, &userAvatars[10].Bot,
		&userAvatars[11].UserId, &userAvatars[11].Url, &userAvatars[11].Name, &userAvatars[11].Bot,
		&userAvatars[12].UserId, &userAvatars[12].Url, &userAvatars[12].Name, &userAvatars[12].Bot,
		&userAvatars[13].UserId, &userAvatars[13].Url, &userAvatars[13].Name, &userAvatars[13].Bot,
		&userAvatars[14].UserId, &userAvatars[14].Url, &userAvatars[14].Name, &userAvatars[14].Bot,
		&userAvatars[15].UserId, &userAvatars[15].Url, &userAvatars[15].Name, &userAvatars[15].Bot,
		&userAvatars[16].UserId, &userAvatars[16].Url, &userAvatars[16].Name, &userAvatars[16].Bot,
		&userAvatars[17].UserId, &userAvatars[17].Url, &userAvatars[17].Name, &userAvatars[17].Bot,
		&userAvatars[18].UserId, &userAvatars[18].Url, &userAvatars[18].Name, &userAvatars[18].Bot,
		&userAvatars[19].UserId, &userAvatars[19].Url, &userAvatars[19].Name, &userAvatars[19].Bot,
		&userAvatars[20].UserId, &userAvatars[20].Url, &userAvatars[20].Name, &userAvatars[20].Bot,
		&userAvatars[21].UserId, &userAvatars[21].Url, &userAvatars[21].Name, &userAvatars[21].Bot,
		&userAvatars[22].UserId, &userAvatars[22].Url, &userAvatars[22].Name, &userAvatars[22].Bot,
		&userAvatars[23].UserId, &userAvatars[23].Url, &userAvatars[23].Name, &userAvatars[23].Bot,
		&userAvatars[24].UserId, &userAvatars[24].Url, &userAvatars[24].Name, &userAvatars[24].Bot,
		&userAvatars[25].UserId, &userAvatars[25].Url, &userAvatars[25].Name, &userAvatars[25].Bot,
		&userAvatars[26].UserId, &userAvatars[26].Url, &userAvatars[26].Name, &userAvatars[26].Bot,
		&userAvatars[27].UserId, &userAvatars[27].Url, &userAvatars[27].Name, &userAvatars[27].Bot,
		&userAvatars[28].UserId, &userAvatars[28].Url, &userAvatars[28].Name, &userAvatars[28].Bot,
		&userAvatars[29].UserId, &userAvatars[29].Url, &userAvatars[29].Name, &userAvatars[29].Bot,
		&userAvatars[30].UserId, &userAvatars[30].Url, &userAvatars[30].Name, &userAvatars[30].Bot,
		&userAvatars[31].UserId, &userAvatars[31].Url, &userAvatars[31].Name, &userAvatars[31].Bot,
	)

	if err != nil {
		return err
	}

	item.LastEventAt = time.UnixMilli(lastEventMilli)

	item.CurrentMinimapIndicator = compareWeeks(week, item.LastEventAt)

	item.LastEventTimestamp = item.LastEventAt.Unix()

	return nil
}

const mrListInProgressCondition = `occured_on.week = ?
	AND events.merge_request_event_type IN (9, 15)
	AND metrics.approved = 0
	AND metrics.merged = 0
	AND metrics.closed = 0
	AND author.bot = 0
	AND user.bot = 0`

func (s *Store) GetMergeRequestInProgressList(date time.Time, teamMembers []int64, nullUserId int64) ([]MergeRequestListItemData, error) {
	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	db, err := sql.Open(s.DriverName, s.DbUrl)

	if err != nil {
		return nil, err
	}

	week := util.GetFormattedWeek(date)

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE %s
			%s
		GROUP BY mr.id;`,
		mrListDataSelect, mrListTables, mrListInProgressCondition, usersInTeamConditionQuery)

	defer db.Close()

	queryParams := make([]interface{}, len(teamMembers)+2)
	queryParams[0] = week
	for i, v := range teamMembers {
		queryParams[i+1] = v
	}

	rows, err := db.QueryContext(s.Context, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var mergeRequests []MergeRequestListItemData
	var userAvatars = make([]ListUserInfo, 2+3*10)

	for rows.Next() {
		var item MergeRequestListItemData

		if err := scanMergeRequestListItemRow(&item, userAvatars, rows, week); err != nil {
			return nil, err
		}

		uniqueUsersMap := make(map[int64]bool)
		for _, userAvatar := range userAvatars {
			if len(item.Actors) >= iMAX_USER_AVATARS_LEN {
				break
			}

			if userAvatar.UserId != nullUserId && !userAvatar.Bot && !uniqueUsersMap[userAvatar.UserId] {
				uniqueUsersMap[userAvatar.UserId] = true
				item.Actors = append(item.Actors, userAvatar)
			}
		}

		mergeRequests = append(mergeRequests, item)
	}

	return mergeRequests, nil
}

const mrListReadyToMergeCondition = `
	events.merge_request_event_type IN (2, 7, 9, 11, 15)
AND metrics.approved = 1
	AND metrics.merged = 0
	AND metrics.closed = 0
	AND author.bot = 0
	AND user.bot = 0`

func (s *Store) GetMergeRequestReadyToMergeList(teamMembers []int64, nullUserId int64, week string) ([]MergeRequestListItemData, error) {
	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	db, err := sql.Open(s.DriverName, s.DbUrl)

	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE %s
			%s
		GROUP BY mr.id
		ORDER BY
			last_updated_at.year ASC,
			last_updated_at.month ASC,
			last_updated_at.day ASC;`,
		mrListDataSelect, mrListTables, mrListReadyToMergeCondition, usersInTeamConditionQuery)

	defer db.Close()

	queryParams := make([]interface{}, len(teamMembers))
	for i, v := range teamMembers {
		queryParams[i] = v
	}

	rows, err := db.QueryContext(s.Context, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var mergeRequests []MergeRequestListItemData
	var userAvatars = make([]ListUserInfo, 2+3*10)

	for rows.Next() {
		var item MergeRequestListItemData

		if err := scanMergeRequestListItemRow(&item, userAvatars, rows, week); err != nil {
			return nil, err
		}

		uniqueUsersMap := make(map[int64]bool)
		for _, userAvatar := range userAvatars {
			if len(item.Actors) >= iMAX_USER_AVATARS_LEN {
				break
			}

			if userAvatar.UserId != nullUserId && !userAvatar.Bot && !uniqueUsersMap[userAvatar.UserId] {
				uniqueUsersMap[userAvatar.UserId] = true
				item.Actors = append(item.Actors, userAvatar)
			}
		}

		mergeRequests = append(mergeRequests, item)
	}

	return mergeRequests, nil
}

const mrListWaitingForReviewCondition = `
	events.merge_request_event_type IN (2, 7, 9, 11, 15)
	AND	metrics.reviewed = 0
		AND metrics.approved = 0
  	AND metrics.merged = 0
  	AND metrics.closed = 0
		AND author.bot = 0
		AND user.bot = 0
		AND mr.id NOT IN (
			SELECT DISTINCT events.merge_request
			FROM transform_merge_request_events AS events
			JOIN transform_dates AS occured_on ON occured_on.id = events.occured_on
			WHERE occured_on.week = ?
			AND events.merge_request_event_type = 9
		)`

func (s *Store) GetMergeRequestWaitingForReviewList(teamMembers []int64, date time.Time, nullUserId int64) ([]MergeRequestListItemData, error) {
	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	db, err := sql.Open(s.DriverName, s.DbUrl)

	if err != nil {
		return nil, err
	}

	week := util.GetFormattedWeek(date)

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE %s
			%s
		GROUP BY mr.id
		ORDER BY
			last_updated_at.year ASC,
			last_updated_at.month ASC,
			last_updated_at.day ASC;`,
		mrListDataSelect, mrListTables, mrListWaitingForReviewCondition, usersInTeamConditionQuery)

	defer db.Close()

	queryParams := make([]interface{}, len(teamMembers)+1)
	queryParams[0] = week
	for i, v := range teamMembers {
		queryParams[i+1] = v
	}

	rows, err := db.QueryContext(s.Context, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var mergeRequests []MergeRequestListItemData
	var userAvatars = make([]ListUserInfo, 2+3*10)

	for rows.Next() {
		var item MergeRequestListItemData

		if err := scanMergeRequestListItemRow(&item, userAvatars, rows, week); err != nil {
			return nil, err
		}

		uniqueUsersMap := make(map[int64]bool)
		for _, userAvatar := range userAvatars {
			if len(item.Actors) >= iMAX_USER_AVATARS_LEN {
				break
			}

			if userAvatar.UserId != nullUserId && !userAvatar.Bot && !uniqueUsersMap[userAvatar.UserId] {
				uniqueUsersMap[userAvatar.UserId] = true
				item.Actors = append(item.Actors, userAvatar)
			}
		}

		mergeRequests = append(mergeRequests, item)
	}

	return mergeRequests, nil
}

const mrListMergedCondition = `occured_on.week = ?
AND events.merge_request_event_type = 11
AND metrics.merged = 1
AND metrics.closed = 1
AND author.bot = 0
AND user.bot = 0`

func (s *Store) GetMergeRequestMergedList(date time.Time, teamMembers []int64, nullUserId int64) ([]MergeRequestListItemData, error) {
	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	week := util.GetFormattedWeek(date)

	db, err := sql.Open(s.DriverName, s.DbUrl)

	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE %s
			%s
		GROUP BY mr.id
		ORDER BY
			last_updated_at.year DESC,
			last_updated_at.month DESC,
			last_updated_at.day DESC;`,
		mrListDataSelect, mrListTables, mrListMergedCondition, usersInTeamConditionQuery)

	defer db.Close()

	queryParams := make([]interface{}, len(teamMembers)+1)
	queryParams[0] = week
	for i, v := range teamMembers {
		queryParams[i+1] = v
	}

	rows, err := db.QueryContext(s.Context, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var mergeRequests []MergeRequestListItemData
	var userAvatars = make([]ListUserInfo, 2+3*10)

	for rows.Next() {
		var item MergeRequestListItemData

		if err := scanMergeRequestListItemRow(&item, userAvatars, rows, week); err != nil {
			return nil, err
		}

		uniqueUsersMap := make(map[int64]bool)
		for _, userAvatar := range userAvatars {
			if len(item.Actors) >= iMAX_USER_AVATARS_LEN {
				break
			}

			if userAvatar.UserId != nullUserId && !userAvatar.Bot && !uniqueUsersMap[userAvatar.UserId] {
				uniqueUsersMap[userAvatar.UserId] = true
				item.Actors = append(item.Actors, userAvatar)
			}
		}

		mergeRequests = append(mergeRequests, item)
	}

	return mergeRequests, nil
}

const mrListClosedCondition = `occured_on.week = ?
AND metrics.merged = 0
AND metrics.closed = 1
AND author.bot = 0
AND user.bot = 0
AND events.timestamp = (
  	SELECT MAX(e2.timestamp)
    FROM transform_merge_request_events e2
    WHERE e2.merge_request = events.merge_request
	AND e2.merge_request_event_type = 7
)`

func (s *Store) GetMergeRequestClosedList(date time.Time, teamMembers []int64, nullUserId int64) ([]MergeRequestListItemData, error) {
	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	week := util.GetFormattedWeek(date)

	db, err := sql.Open(s.DriverName, s.DbUrl)

	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE %s
			%s
		GROUP BY mr.id
		ORDER BY
			last_updated_at.year DESC,
			last_updated_at.month DESC,
			last_updated_at.day DESC;`,
		mrListDataSelect, mrListTables, mrListClosedCondition, usersInTeamConditionQuery)

	defer db.Close()

	queryParams := make([]interface{}, len(teamMembers)+1)
	queryParams[0] = week
	for i, v := range teamMembers {
		queryParams[i+1] = v
	}

	rows, err := db.QueryContext(s.Context, query, queryParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var mergeRequests []MergeRequestListItemData
	var userAvatars = make([]ListUserInfo, 2+3*10)

	for rows.Next() {
		var item MergeRequestListItemData

		if err := scanMergeRequestListItemRow(&item, userAvatars, rows, week); err != nil {
			return nil, err
		}

		uniqueUsersMap := make(map[int64]bool)
		for _, userAvatar := range userAvatars {
			if len(item.Actors) >= iMAX_USER_AVATARS_LEN {
				break
			}

			if userAvatar.UserId != nullUserId && !userAvatar.Bot && !uniqueUsersMap[userAvatar.UserId] {
				uniqueUsersMap[userAvatar.UserId] = true
				item.Actors = append(item.Actors, userAvatar)
			}
		}

		mergeRequests = append(mergeRequests, item)
	}

	return mergeRequests, nil
}

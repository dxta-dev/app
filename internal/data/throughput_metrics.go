package data

import (
	"database/sql"
	"fmt"
	"strings"
)

type CodeChangesCount struct {
	Count int
	Week  string
}

func (s *Store) GetTotalCodeChanges(weeks []string, teamMembers []int64) (map[string]CodeChangesCount, float64, error) {

	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	query := fmt.Sprintf(`
	SELECT
		SUM(metrics.mr_size) AS total_mr_size,
		dates.week
	FROM transform_merge_request_metrics AS metrics
	JOIN transform_merge_request_fact_dates_junk AS dates_junk
	ON metrics.dates_junk = dates_junk.id
	JOIN transform_dates AS dates
	ON dates_junk.merged_at = dates.id
	JOIN transform_merge_request_fact_users_junk AS uj
	ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
	ON uj.author = author.id
	WHERE dates.week IN (%s)
	AND author.bot = 0
	%s
	GROUP BY dates.week;`,
		placeholders,
		usersInTeamConditionQuery)

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, 0, err
	}

	defer db.Close()

	queryParams := make([]interface{}, len(weeks)+len(teamMembers))
	for i, v := range weeks {
		queryParams[i] = v
	}
	for i, v := range teamMembers {
		queryParams[i+len(weeks)] = v
	}

	rows, err := db.Query(query, queryParams...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	codeChangesByWeek := make(map[string]CodeChangesCount)

	for rows.Next() {
		var codeChangesCount CodeChangesCount

		if err = rows.Scan(&codeChangesCount.Count, &codeChangesCount.Week); err != nil {
			return nil, 0, err
		}
		codeChangesByWeek[codeChangesCount.Week] = codeChangesCount
	}

	totalCodeChangesCount := 0
	numOfWeeksWithCodeChanges := len(codeChangesByWeek)

	for _, week := range weeks {
		if codeChangesByWeek[week].Count >= 0 {
			totalCodeChangesCount += codeChangesByWeek[week].Count
		}

		if _, ok := codeChangesByWeek[week]; !ok {
			codeChangesByWeek[week] = CodeChangesCount{
				Count: 0,
				Week:  week,
			}
		}
	}

	averageCodeChangesByXWeeks := float64(totalCodeChangesCount) / float64(numOfWeeksWithCodeChanges)

	return codeChangesByWeek, averageCodeChangesByXWeeks, nil
}

type CommitCountByWeek struct {
	Week  string
	Count int
}

func (s *Store) GetTotalCommits(weeks []string, teamMembers []int64) (map[string]CommitCountByWeek, float64, error) {

	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND actor.external_id IN (%s)", teamMembersPlaceholders)
	}

	query := fmt.Sprintf(`
		SELECT
			COUNT (*),
			commitedAt.week
		FROM transform_merge_request_events AS ev
		JOIN transform_dates AS commitedAt
		ON commitedAt.id = ev.commited_at
		JOIN transform_forge_users AS actor
		ON ev.actor = actor.id
		WHERE ev.merge_request_event_type = 9
		AND commitedAt.week IN (%s)
		AND actor.bot = 0
		%s
		GROUP BY commitedAt.week;`,
		placeholders,
		usersInTeamConditionQuery)

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, 0, err
	}

	defer db.Close()

	queryParams := make([]interface{}, len(weeks)+len(teamMembers))
	for i, v := range weeks {
		queryParams[i] = v
	}
	for i, v := range teamMembers {
		queryParams[i+len(weeks)] = v
	}

	rows, err := db.Query(query, queryParams...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	commitCountByWeeks := make(map[string]CommitCountByWeek)

	for rows.Next() {
		var commitCount CommitCountByWeek

		if err := rows.Scan(&commitCount.Count, &commitCount.Week); err != nil {
			return nil, 0, err
		}
		commitCountByWeeks[commitCount.Week] = commitCount
	}

	totalCommitCount := 0
	numOfWeeksWithCommits := len(commitCountByWeeks)

	for _, week := range weeks {
		totalCommitCount += commitCountByWeeks[week].Count
		if _, ok := commitCountByWeeks[week]; !ok {
			commitCountByWeeks[week] = CommitCountByWeek{
				Week:  week,
				Count: 0,
			}
		}
	}

	averageCommitCountByXWeeks := float64(totalCommitCount) / float64(numOfWeeksWithCommits)

	return commitCountByWeeks, averageCommitCountByXWeeks, nil
}

func (s *Store) GetTotalMrsOpened(weeks []string, teamMembers []int64) (map[string]MrCountByWeek, float64, error) {

	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	query := fmt.Sprintf(`
	SELECT
		COUNT (*),
		opened_dates.week
	FROM transform_merge_request_metrics AS metrics
	JOIN transform_merge_request_fact_dates_junk AS dates_junk
	ON metrics.dates_junk = dates_junk.id
	JOIN transform_dates AS opened_dates
	ON dates_junk.opened_at = opened_dates.id
	JOIN transform_merge_request_fact_users_junk AS uj
	ON metrics.users_junk = uj.id
	JOIN transform_forge_users AS author
	ON uj.author = author.id
	WHERE opened_dates.week IN (%s)
	AND author.bot = 0
	%s
	GROUP BY opened_dates.week`,
		placeholders,
		usersInTeamConditionQuery)

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, 0, err
	}

	defer db.Close()

	queryParams := make([]interface{}, len(weeks)+len(teamMembers))
	for i, v := range weeks {
		queryParams[i] = v
	}
	for i, v := range teamMembers {
		queryParams[i+len(weeks)] = v
	}

	rows, err := db.Query(query, queryParams...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	mrCountByWeeks := make(map[string]MrCountByWeek)

	for rows.Next() {
		var prCount MrCountByWeek

		if err := rows.Scan(&prCount.Count, &prCount.Week); err != nil {
			return nil, 0, err
		}
		mrCountByWeeks[prCount.Week] = prCount
	}

	totalMRCount := 0
	numOfWeeksWithMR := len(mrCountByWeeks)

	for _, week := range weeks {
		totalMRCount += mrCountByWeeks[week].Count
		if _, ok := mrCountByWeeks[week]; !ok {
			mrCountByWeeks[week] = MrCountByWeek{
				Week:  week,
				Count: 0,
			}
		}
	}

	averageMRCountByXWeeks := float64(totalMRCount) / float64(numOfWeeksWithMR)

	return mrCountByWeeks, averageMRCountByXWeeks, nil
}

type TotalReviewsByWeek struct {
	Week  string
	Count int
}

func (s *Store) GetTotalReviews(weeks []string, teamMembers []int64) (map[string]TotalReviewsByWeek, float64, error) {

	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND actor.external_id IN (%s)", teamMembersPlaceholders)
	}

	query := fmt.Sprintf(`
		SELECT
			COUNT (*),
			occuredAt.week
		FROM transform_merge_request_events AS ev
		JOIN transform_dates AS occuredAt
		ON occuredAt.id = ev.occured_on
		JOIN transform_forge_users AS actor
		ON ev.actor = actor.id
		WHERE ev.merge_request_event_type = 15
		AND occuredAt.week IN (%s)
		AND actor.bot = 0
		%s
		GROUP BY occuredAt.week;`,
		placeholders,
		usersInTeamConditionQuery)

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, 0, err
	}

	defer db.Close()

	queryParams := make([]interface{}, len(weeks)+len(teamMembers))
	for i, v := range weeks {
		queryParams[i] = v
	}
	for i, v := range teamMembers {
		queryParams[i+len(weeks)] = v
	}

	rows, err := db.Query(query, queryParams...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	totalReviewsByWeek := make(map[string]TotalReviewsByWeek)

	for rows.Next() {
		var reviewCount TotalReviewsByWeek

		if err := rows.Scan(&reviewCount.Count, &reviewCount.Week); err != nil {
			return nil, 0, err
		}
		totalReviewsByWeek[reviewCount.Week] = reviewCount
	}

	totalReviewsCount := 0
	numOfWeeksWithReviews := len(totalReviewsByWeek)

	for _, week := range weeks {
		totalReviewsCount += totalReviewsByWeek[week].Count
		if _, ok := totalReviewsByWeek[week]; !ok {
			totalReviewsByWeek[week] = TotalReviewsByWeek{
				Week:  week,
				Count: 0,
			}
		}
	}

	averageReviewsByXWeeks := float64(totalReviewsCount) / float64(numOfWeeksWithReviews)

	return totalReviewsByWeek, averageReviewsByXWeeks, nil
}

type MergeFrequencyByWeek struct {
	Week   string
	Amount float32
}

func (s *Store) GetMergeFrequency(weeks []string, teamMembers []int64) (map[string]MergeFrequencyByWeek, float64, error) {
	placeholders := strings.Repeat("?,", len(weeks)-1) + "?"

	usersInTeamConditionQuery := ""
	if len(teamMembers) > 0 {
		teamMembersPlaceholders := strings.Repeat("?,", len(teamMembers)-1) + "?"
		usersInTeamConditionQuery = fmt.Sprintf("AND author.external_id IN (%s)", teamMembersPlaceholders)
	}

	query := fmt.Sprintf(`
		SELECT
			CAST(COUNT (*) AS REAL) / 7,
			merged_dates.week
		FROM transform_merge_request_metrics AS metrics
		JOIN transform_merge_request_fact_dates_junk AS dates_junk
		ON metrics.dates_junk = dates_junk.id
		JOIN transform_dates AS merged_dates
		ON dates_junk.merged_at = merged_dates.id
		JOIN transform_merge_request_fact_users_junk AS uj
		ON metrics.users_junk = uj.id
		JOIN transform_forge_users AS author
		ON uj.author = author.id
		WHERE merged_dates.week IN (%s)
		AND author.bot = 0
		%s
		GROUP BY merged_dates.week`,
		placeholders,
		usersInTeamConditionQuery)

	db, err := sql.Open("libsql", s.DbUrl)

	if err != nil {
		return nil, 0, err
	}

	defer db.Close()

	queryParams := make([]interface{}, len(weeks)+len(teamMembers))
	for i, v := range weeks {
		queryParams[i] = v
	}
	for i, v := range teamMembers {
		queryParams[i+len(weeks)] = v
	}

	rows, err := db.Query(query, queryParams...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	mergeFrequencyByWeek := make(map[string]MergeFrequencyByWeek)

	for rows.Next() {
		var mergeFreq MergeFrequencyByWeek

		if err := rows.Scan(&mergeFreq.Amount, &mergeFreq.Week); err != nil {
			return nil, 0, err
		}
		mergeFrequencyByWeek[mergeFreq.Week] = mergeFreq
	}

	totalMergeFrequencyCount := 0.0
	numOfWeeksWithMergeFrequency := len(mergeFrequencyByWeek)

	for _, week := range weeks {
		totalMergeFrequencyCount += float64(mergeFrequencyByWeek[week].Amount)
		if _, ok := mergeFrequencyByWeek[week]; !ok {
			mergeFrequencyByWeek[week] = MergeFrequencyByWeek{
				Week:   week,
				Amount: 0,
			}
		}
	}

	averageMergeFrequencyByXWeeks := totalMergeFrequencyCount / float64(numOfWeeksWithMergeFrequency)

	return mergeFrequencyByWeek, averageMergeFrequencyByXWeeks, nil
}

func (s *Store) GetDeployFrequency(weeks []string) (interface{}, error) {
	return nil, nil
}

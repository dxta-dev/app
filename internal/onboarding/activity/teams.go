package activity

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dxta-dev/app/internal/onboarding"
)

func (ta *TenantActivities) CreateTeamMembers(ctx context.Context,
	DBURL string,
	members MembersRecordMap,
	organizationID int64) ([]MemberRecord, error) {
	db, err := onboarding.GetCachedTenantDB(ta.DBConnections, DBURL, ctx)

	if err != nil {
		return nil, errors.New("failed to get cached tenant db to upsert teams: " + err.Error())
	}

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return nil, errors.New("failed to begin transaction to upsert teams: " + err.Error())
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	args := make([]any, 0)
	values := make([]string, 0)
	idsToUpdate := make([]string, 0)

	for _, member := range members {
		args = append(args, member.Name, member.Email, member.Login)
		values = append(values, "(?, ?, ?)")
		idsToUpdate = append(idsToUpdate, fmt.Sprintf("%d", *member.GithubMemberId))
	}

	query := fmt.Sprintf(`
		INSERT INTO members
			(name, email, username)
		VALUES
			%s
		RETURNING id, username;
	`, strings.Join(values, ", "))

	rows, err := tx.QueryContext(ctx, query,
		args...)

	if err != nil {
		return nil, errors.New("failed to create members: " + err.Error())
	}

	caseStatements := make([]string, 0)
	newMembers := make([]MemberRecord, 0)

	for rows.Next() {
		var res struct {
			ID       *int64
			Username *string
		}

		if err := rows.Scan(&res.ID, &res.Username); err != nil {
			return nil, errors.New("failed to scan create member result: " + err.Error())
		}

		member, ok := members[*res.Username]

		if !ok {
			return nil, errors.New("failed to get a member record from map")
		}

		member.MemberID = res.ID
		newMembers = append(newMembers, member)

		caseStatements = append(caseStatements, fmt.Sprintf("WHEN username = '%s' THEN %d", *res.Username, *res.ID))

	}
	caseClause := strings.Join(caseStatements, "\n    ")
	whereInClause := strings.Join(idsToUpdate, ", ")

	query = fmt.Sprintf(`
		UPDATE 
			github_members 
		SET member_id = CASE 
			%s 
		ELSE member_id 
		END 
		WHERE id IN (%s)
	`, caseClause, whereInClause)

	_, err = tx.ExecContext(ctx, query)

	if err != nil {
		return nil, errors.New("failed to update github members reference: " + err.Error())
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return newMembers, nil
}

func (ta *TenantActivities) JoinTeamsMembers(
	ctx context.Context,
	DBURL string,
	newMembers []MemberRecord,
) (bool, error) {
	db, err := onboarding.GetCachedTenantDB(ta.DBConnections, DBURL, ctx)

	if err != nil {
		return false, errors.New("failed to get cached tenant db to upsert teams: " + err.Error())
	}

	args := make([]any, 0)
	values := make([]string, 0)

	for _, member := range newMembers {
		for _, team := range member.Teams {
			args = append(args, team.TeamID, member.MemberID)
			values = append(values, "(?, ?)")
		}

	}

	query := fmt.Sprintf(`
		INSERT INTO teams__members
			(team_id, member_id)
		VALUES %s ;`, strings.Join(values, ", "))

	_, err = db.ExecContext(ctx, query, args...)

	if err != nil {
		return false, errors.New("failed to insert into teams_members: " + err.Error())
	}

	return true, nil
}

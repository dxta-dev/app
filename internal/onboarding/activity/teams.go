package activity

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dxta-dev/app/internal/onboarding"
)

func (ta *TenantActivities) CreateTeamMember(
	ctx context.Context,
	DBURL string,
	member MemberRecord,
	organizationID int64,
) (*MemberRecord, error) {
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

	row := db.QueryRowContext(ctx, `
		INSERT INTO members
			(name, email)
		VALUES
			(?, ?)
		RETURNING id;`,
		member.Name, member.Email)

	var memberId int64

	member.MemberID = &memberId

	err = row.Scan(&memberId)

	if err != nil {
		return nil, errors.New("Issue creating member: " + err.Error())
	}

	_, err = tx.Exec(`
		UPDATE 
			github_members 
		SET 
			member_id = ? 
		WHERE id = ?`,
		memberId, member.GithubMemberId)

	if err != nil {
		return nil, errors.New("Issue while updating member_id in github member: " + err.Error())
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &member, err
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
			args = append(args, []any{team.TeamID, member.MemberID}...)
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

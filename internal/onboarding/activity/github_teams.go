package activity

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/dxta-dev/app/internal/onboarding"
)

type GithubActivities struct {
	ClientMap    *sync.Map
	GithubConfig onboarding.GithubConfig
}

func NewGithubActivities(clientMap *sync.Map, githubConfig onboarding.GithubConfig) *GithubActivities {
	return &GithubActivities{
		ClientMap:    clientMap,
		GithubConfig: githubConfig,
	}
}

func (ga *GithubActivities) GetTeams(ctx context.Context, installationId int64, organizationName string) (onboarding.Teams, error) {

	client, err := onboarding.GetCachedGithubInstallationClient(ga.ClientMap, installationId, ga.GithubConfig)

	if err != nil {
		return nil, errors.New("could not create new installation client to get github teams " + err.Error())
	}

	return client.GetTeams(ctx, organizationName, &onboarding.ListOptions{PerPage: 100, MaxPages: 50})
}

func (ga *GithubActivities) GetTeamMembers(ctx context.Context, installationId int64, organizationName string, teamSlug string) (onboarding.Members, error) {
	client, err := onboarding.GetCachedGithubInstallationClient(ga.ClientMap, installationId, ga.GithubConfig)

	if err != nil {
		return nil, errors.New("could not create new installation client to get github members " + err.Error())
	}

	return client.GetTeamMembers(ctx, organizationName, teamSlug, &onboarding.ListOptions{PerPage: 100, MaxPages: 50})
}

func (ga *GithubActivities) GetTeamMemberEmail(ctx context.Context, installationId int64, teamMember onboarding.Member, teamSlug string) (*onboarding.Member, error) {
	client, err := onboarding.GetCachedGithubInstallationClient(ga.ClientMap, installationId, ga.GithubConfig)

	if err != nil {
		return nil, errors.New("could not create new installation client to get github members " + err.Error())
	}

	return client.GetTeamMemberWithEmail(ctx, teamMember)
}

type UpsertTeamRes struct {
	githubTeamId int64
	teamId       int64
}

func upsertTeam(ctx context.Context, tx *sql.Tx, teamIdPtr *int64, teamNamePtr *string, organizationId, githubOrganizationId int64) (*UpsertTeamRes, error) {

	if tx == nil || teamIdPtr == nil || teamNamePtr == nil {
		return nil, errors.New("invalid params in upsertTeams")
	}

	var githubTeamId int64

	rows := tx.QueryRowContext(ctx, `
		SELECT id 
		FROM github_teams 
		WHERE external_id = ?;`, teamIdPtr)

	if err := rows.Scan(&githubTeamId); err != nil {
		if err != sql.ErrNoRows {
			return nil, errors.New("failed to retrieve github_team: " + err.Error())
		}
	}

	var teamId int64

	if githubTeamId == 0 {
		rows = tx.QueryRowContext(ctx, `
		INSERT INTO github_teams 
            (name, external_id, github_organization_id) 
        VALUES 
            (?, ?, ?) 
        RETURNING 
            id;`,
			teamNamePtr, teamIdPtr, githubOrganizationId)

		if err := rows.Scan(&githubTeamId); err != nil {
			return nil, errors.New("failed to upsert github_teams: " + err.Error())
		}

		rows = tx.QueryRowContext(ctx, `
		INSERT INTO teams 
            (name, organization_id) 
        VALUES 
            (?, ?) 
        RETURNING 
            id;`,
			teamNamePtr, organizationId)

		if err := rows.Scan(&teamId); err != nil {
			return nil, errors.New("failed to upsert teams: " + err.Error())
		}
	}

	return &UpsertTeamRes{githubTeamId, teamId}, nil
}

func upsertMember(ctx context.Context, tx *sql.Tx, teamId int64, githubTeamId int64, member onboarding.Member) (bool, error) {
	rowRes := tx.QueryRowContext(ctx, `
			INSERT INTO github_members
				(external_id, username, email)
			VALUES
				(?, ?, ?)
			ON CONFLICT 
				(external_id) 
			DO UPDATE SET 
				username = excluded.username, 
				email = excluded.email 
			RETURNING id, member_id`,
		member.ID, member.Login, member.Email)

	var githubMemberId int64
	var memberRefId *int64

	if err := rowRes.Scan(&githubMemberId, &memberRefId); err != nil {
		return false, errors.New("failed to upsert github_member: " + err.Error())
	}

	_, err := tx.Exec(`
			INSERT OR IGNORE INTO github_teams__github_members
				(github_team_id, github_member_id)
			VALUES
				(?, ?);`,
		githubTeamId, githubMemberId)

	if err != nil {
		return false, errors.New("Issue creating github_teams__github_members: " + err.Error())
	}

	if memberRefId == nil {
		name := member.Name

		if name == nil {
			defaultName := "DXTA member"
			name = &defaultName
		}

		rowRes = tx.QueryRowContext(ctx, `
			INSERT INTO members
				(name, email)
			VALUES
				(?, ?)
			RETURNING id;`,
			name, member.Email)

		var memberId int64

		err = rowRes.Scan(&memberId)

		if err != nil {
			return false, errors.New("Issue creating member: " + err.Error())
		}

		_, err = tx.Exec(`
				UPDATE 
					github_members 
				SET 
					member_id = ? 
				WHERE id = ?`,
			memberId, githubMemberId)

		if err != nil {
			return false, errors.New("Issue while updating member_id in github member: " + err.Error())
		}

		_, err = tx.Exec(`
			INSERT INTO teams__members
				(team_id, member_id)
			VALUES
				(?, ?);`,
			teamId, memberId)

		if err != nil {
			return false, errors.New("Issue creating teams__members: " + err.Error())
		}
	}

	if memberRefId != nil && teamId != 0 {
		_, err = tx.Exec(`
			INSERT INTO teams__members
				(team_id, member_id)
			VALUES
				(?, ?);`,
			teamId, memberRefId)

		if err != nil {
			return false, errors.New("Issue creating teams__members: " + err.Error())
		}
	}

	return true, nil
}

func (ta *TenantActivities) UpsertTeamAndMembers(ctx context.Context, team onboarding.Team, members onboarding.Members, DBURL string, githubOrganizationId int64, organizationId int64) (res bool, err error) {
	db, err := onboarding.GetCachedTenantDB(ta.DBConnections, DBURL, ctx)

	if err != nil {
		return false, errors.New("failed to get cached tenant db to upsert teams and members: " + err.Error())
	}

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return false, errors.New("failed to begin transaction to upsert teams and members: " + err.Error())
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	teamData, err := upsertTeam(ctx, tx, team.ID, team.Name, organizationId, githubOrganizationId)

	if err != nil {
		return false, errors.New("failed to upsert teams: " + err.Error())
	}

	for _, member := range members {
		_, err := upsertMember(ctx, tx, teamData.teamId, teamData.githubTeamId, member)

		if err != nil {
			return false, errors.New("failed to upsert member: " + err.Error())
		}
	}

	return true, nil
}

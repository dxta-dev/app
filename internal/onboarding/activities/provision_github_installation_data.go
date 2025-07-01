package activities

import (
	"context"
	"database/sql"
	"fmt"

	internal_api_data "github.com/dxta-dev/app/internal/internal-api/data"
	"github.com/dxta-dev/app/internal/onboarding/data"
	"github.com/google/go-github/v72/github"
)

func (activity *GithubActivities) GetGithubInstallation(
	ctx context.Context,
	installationId int64,
) (*github.Installation, error) {
	return data.GetGithubInstallation(installationId, activity.githubConfig.GithubAppClient, ctx)
}

func (activity *DBActivities) SyncGithubInstallationDataToTenant(ctx context.Context, installationId int64,
	installationOrgName string,
	installationOrgId int64,
	organizationId string,
	dbUrl string) (*data.SyncGithubDataResult, error) {
	cacheKey := dbUrl
	db, ok := activity.connections.Load(cacheKey)

	if !ok {
		tenantDB, err := internal_api_data.NewTenantDB(cacheKey, ctx)
		db = tenantDB.DB

		if err != nil {
			return nil, err
		}

		activity.connections.Store(cacheKey, db)
	}

	res, err := data.SyncGithubInstallationDataToTenant(
		installationId,
		installationOrgName,
		installationOrgId,
		organizationId,
		db.(*sql.DB), ctx,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (activity *GithubActivities) GetInstallationTeams(
	ctx context.Context,
	installationOrgName string,
	installationId int64,
) ([]*github.Team, error) {
	client, err := activity.newInstallationClient(installationId)

	if err != nil {
		fmt.Printf("Could not create new installation client. Error: %v", err.Error())
		return nil, err
	}

	return data.GetInstallationTeams(ctx, installationOrgName, client)
}

func (activity *GithubActivities) GetInstallationTeamMembers(
	ctx context.Context,
	installationId int64,
	installationOrgName string,
	teamSlug string,
) ([]*github.User, error) {
	client, err := activity.newInstallationClient(installationId)

	if err != nil {
		fmt.Printf("Could not create new installation client. Error: %v", err.Error())
		return nil, err
	}
	return data.GetInstallationTeamMembers(ctx, installationOrgName, teamSlug, client)
}

func (activity *GithubActivities) GetInstallationTeamMembersWithEmails(ctx context.Context, installationId int64, members []*github.User) (data.ExtendedMembers, error) {

	client, err := activity.newInstallationClient(installationId)

	if err != nil {
		fmt.Printf("Could not create new installation client. Error: %v", err.Error())
		return nil, err
	}

	return data.GetInstallationTeamMembersWithEmails(ctx, members, client)
}

func (activity *DBActivities) SyncTeamsAndMembersToTenant(
	ctx context.Context,
	teamWithMembers data.TeamWithMembers,
	dbUrl string,
	githubOrganizationId int64,
	organizationId int64,
) (bool, error) {
	cacheKey := dbUrl
	d, ok := activity.connections.Load(cacheKey)

	if !ok {
		tenantDB, err := internal_api_data.NewTenantDB(cacheKey, ctx)
		d = tenantDB.DB

		if err != nil {
			return false, err
		}

		activity.connections.Store(cacheKey, d)
	}
	db := d.(*sql.DB)

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return false, err
	}

	rows := tx.QueryRowContext(ctx, `
		INSERT INTO teams 
            (name, organization_id) 
        VALUES 
            (?, ?) 
        RETURNING 
            id;`,
		teamWithMembers.Team.Name, organizationId)

	var teamId int64

	err = rows.Scan(&teamId)

	if err != nil {
		fmt.Println("Issue creating team")
		_ = tx.Rollback()
		return false, err
	}

	rows = tx.QueryRowContext(ctx, `
		INSERT INTO github_teams 
            (name, external_id, github_organization_id) 
        VALUES 
            (?, ?, ?) 
        RETURNING 
            id;`,
		teamWithMembers.Team.Name, teamWithMembers.Team.ID, githubOrganizationId)

	var githubTeamId int64

	err = rows.Scan(&githubTeamId)

	if err != nil {
		fmt.Println("Issue creating github team")

		_ = tx.Rollback()
		return false, err
	}

	for _, member := range teamWithMembers.Members {
		name := member.Name

		if name == nil {
			defaultName := "DXTA member"
			name = &defaultName
		}

		rowRes := tx.QueryRowContext(ctx, `
			INSERT INTO members
				(name, email)
			VALUES
				(?, ?)
			RETURNING id;`,
			name, member.Email)

		var memberId int64

		err = rowRes.Scan(&memberId)

		if err != nil {
			fmt.Println("Issue creating member")

			_ = tx.Rollback()
			return false, err
		}

		_, err = tx.Exec(`
			INSERT INTO teams__members
				(team_id, member_id)
			VALUES
				(?, ?);`,
			teamId, memberId)

		if err != nil {
			fmt.Println("Issue creating teams__members")
			_ = tx.Rollback()
			return false, err
		}

		rowRes = tx.QueryRowContext(ctx, `
			INSERT INTO github_members
				(external_id, username, email, member_id)
			VALUES
				(?, ?, ?, ?)
			RETURNING id;`,
			member.ID, member.Login, member.Email, memberId)

		var githubMemberId int64

		err = rowRes.Scan(&githubMemberId)

		if err != nil {
			fmt.Println("Issue creating github member")
			_ = tx.Rollback()
			return false, err
		}

		_, err = tx.Exec(`
			INSERT INTO github_teams__github_members
				(github_team_id, github_member_id)
			VALUES
				(?, ?);`,
			githubTeamId, githubMemberId)

		if err != nil {
			fmt.Println("Issue creating github_teams__github_members")
			_ = tx.Rollback()
			return false, err
		}
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

package data

import (
	"context"
	"database/sql"
	"fmt"
)

func SyncTeamsToTenant(tx *sql.Tx, ctx context.Context, teamIdPtr *int64, teamNamePtr *string, organizationId, githubOrganizationId int64) (*struct {
	GithubTeamId int64
	TeamId       int64
}, error) {
	var githubTeamId int64

	rows := tx.QueryRowContext(ctx, `
		INSERT INTO github_teams 
            (name, external_id, github_organization_id) 
        VALUES 
            (?, ?, ?) 
        RETURNING 
            id;`,
		teamNamePtr, teamIdPtr, githubOrganizationId)

	err := rows.Scan(&githubTeamId)

	if err != nil {
		fmt.Println("Issue creating github team")

		_ = tx.Rollback()
		return nil, err
	}

	var teamId int64

	rows = tx.QueryRowContext(ctx, `
		INSERT INTO teams 
            (name, organization_id) 
        VALUES 
            (?, ?) 
        RETURNING 
            id;`,
		teamNamePtr, organizationId)

	err = rows.Scan(&teamId)

	if err != nil {
		fmt.Println("Issue creating team")
		_ = tx.Rollback()
		return nil, err
	}

	return &struct {
		GithubTeamId int64
		TeamId       int64
	}{GithubTeamId: githubTeamId, TeamId: teamId}, nil
}

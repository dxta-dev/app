package data

import (
	"context"
	"database/sql"
)

type SyncGithubDataResult struct {
	OrganizationId       int64
	GithubOrganizationId int64
}

func SyncGithubInstallationDataToTenant(
	installationId int64,
	installationOrgName string,
	installationOrgId int64,
	organizationId string,
	db *sql.DB,
	ctx context.Context,
) (*SyncGithubDataResult, error) {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	rows := tx.QueryRowContext(ctx, `
		INSERT INTO github_organizations
    		(github_app_installation_id, name, external_id)
    	VALUES
			(?, ?, ?) 
		RETURNING id`,
		installationId, installationOrgName, installationOrgId)

	var githubOrganizationId int64

	err = rows.Scan(&githubOrganizationId)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	rows = tx.QueryRowContext(ctx, `
		SELECT 
			id 
		FROM 
			organizations 
		WHERE 
			auth_id = ?;`,
		organizationId)

	var orgId int64
	err = rows.Scan(&orgId)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO 'organizations__github_organizations'
    		('organization_id', 'github_organization_id')
    	VALUES
    	    (?, ?);`,
		orgId, githubOrganizationId)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &SyncGithubDataResult{
		OrganizationId:       orgId,
		GithubOrganizationId: githubOrganizationId,
	}, nil
}

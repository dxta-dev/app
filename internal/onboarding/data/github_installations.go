package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/go-github/v72/github"
)

func GetGithubInstallation(installationId int64, githubAppClient *github.Client, ctx context.Context) (*github.Installation, error) {
	installation, _, err := githubAppClient.Apps.GetInstallation(ctx, installationId)

	if err != nil {
		fmt.Printf("Could not retrieve installation. Error: %v", err.Error())
		return nil, err
	}

	return installation, nil
}

func SyncGithubInstallationDataToTenant(
	installationId int64,
	installationOrgName string,
	installationOrgId int64,
	organizationId string,
	db *sql.DB,
	ctx context.Context,
) error {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
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
		return err
	}

	rows = tx.QueryRowContext(ctx, `
		SELECT id 
		FROM organizations 
		WHERE auth_id = ?;`,
		organizationId)

	var orgId int64
	err = rows.Scan(&orgId)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO 'organizations__github_organizations'
    		('organization_id', 'github_organization_id')
    	VALUES
    	    (?, ?);`,
		orgId, githubOrganizationId)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

package data

import (
	"context"
	"fmt"

	"github.com/google/go-github/v72/github"
)

func (cfg GithubCfg) GetGithubInstallation(installationId int64, ctx context.Context) (*github.Installation, error) {
	githubAppClient := cfg.GithubAppClient

	installation, _, err := githubAppClient.Apps.GetInstallation(ctx, installationId)

	if err != nil {
		fmt.Printf("Could not retrieve installation. Error: %v", err.Error())
		return nil, err
	}

	return installation, nil
}

/* func (d TenantDB) SyncGithubInstallationDataToTenant(
	installationId int64,
	installationOrgName string,
	organizationId string,
	ctx context.Context,
) error {
	tx, err := d.DB.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO github_organizations
    		(github_app_installation_id, name)
    	VALUES
			(?, ?);`,
		installationId, installationOrgName)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO organizations
			(external_id)
    	VALUES
			(?)
    	ON CONFLICT
			(external_id)
    	DO NOTHING;`,
		organizationId)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO 'organizations_github_organizations'
    		('organization_id', 'github_app_installation_id')
    	VALUES
    	    (?, ?);`,
		organizationId, installationId)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
} */

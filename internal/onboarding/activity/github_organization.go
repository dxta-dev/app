package activity

import (
	"context"
	"errors"

	"github.com/dxta-dev/app/internal/onboarding"
)

func (ta *TenantActivities) UpsertGithubOrganization(
	ctx context.Context,
	DBURL string,
	installationId int64,
	installationOrgName string,
	installationOrgId int64,
	organizationId int64,
) (res int64, err error) {
	db, err := onboarding.GetCachedTenantDB(ta.DBConnections, DBURL, ctx)

	if err != nil {
		return 0, err
	}

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	rows := tx.QueryRowContext(ctx, `
		INSERT INTO github_organizations
    		(github_app_installation_id, name, external_id)
    	VALUES
			(?, ?, ?)
		ON CONFLICT (external_id) 
		DO UPDATE SET 
			github_app_installation_id = excluded.github_app_installation_id, 
			name = excluded.name
		RETURNING id`,
		installationId, installationOrgName, installationOrgId)

	var githubOrganizationId int64

	if err = rows.Scan(&githubOrganizationId); err != nil {
		return 0, errors.New("failed to upsert github_organizations: " + err.Error())
	}

	_, err = tx.Exec(`
		INSERT OR IGNORE INTO 'organizations__github_organizations'
    		('organization_id', 'github_organization_id')
    	VALUES
    	    (?, ?)`,
		organizationId, githubOrganizationId)

	if err != nil {
		return 0, errors.New("failed to upsert organizations__github_organizations:" + err.Error())
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.New("failed to commit github organization upsert tx: " + err.Error())
	}

	return githubOrganizationId, nil
}

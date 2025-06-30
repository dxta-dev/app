package activities

import (
	"context"
	"database/sql"

	internal_api_data "github.com/dxta-dev/app/internal/internal-api/data"
	"github.com/dxta-dev/app/internal/onboarding/data"
	"github.com/google/go-github/v72/github"
)

func (activity *GithubActivities) GetGithubInstallation(ctx context.Context, installationId int64) (*github.Installation, error) {
	installations, err := data.GetGithubInstallation(installationId, activity.githubConfig.GithubAppClient, ctx)

	if err != nil {
		return nil, err
	}

	return installations, nil
}

func (activity *DBActivities) SyncGithubInstallationDataToTenant(ctx context.Context, installationId int64,
	installationOrgName string,
	installationOrgId int64,
	organizationId string,
	dbUrl string) (bool, error) {
	cacheKey := dbUrl
	db, ok := activity.connections.Load(cacheKey)

	if !ok {
		tenantDB, err := internal_api_data.NewTenantDB(cacheKey, ctx)
		db = tenantDB.DB

		if err != nil {
			return false, err
		}

		activity.connections.Store(cacheKey, db)
	}

	err := data.SyncGithubInstallationDataToTenant(
		installationId,
		installationOrgName,
		installationOrgId,
		organizationId,
		db.(*sql.DB), ctx,
	)

	if err != nil {
		return false, err
	}

	return true, nil

}

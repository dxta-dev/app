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

	return data.GetInstallationTeamMembers(ctx, installationOrgName, teamSlug, client /* extendWithEmail= */, false)
}

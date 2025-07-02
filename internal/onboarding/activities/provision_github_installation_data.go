package activities

import (
	"context"
	"fmt"

	"github.com/dxta-dev/app/internal/onboarding/data"
	"github.com/google/go-github/v72/github"
)

func (activity *GithubActivities) GetGithubInstallation(
	ctx context.Context,
	installationId int64,
) (*github.Installation, error) {
	return data.GetGithubInstallation(installationId, activity.GithubConfig.GithubAppClient, ctx)
}

func (activity *DBActivities) SyncGithubInstallationDataToTenant(ctx context.Context, installationId int64,
	installationOrgName string,
	installationOrgId int64,
	authId string,
	dbUrl string) (*data.SyncGithubDataResult, error) {
	db, err := activity.GetCachedTenantDB(&activity.Connections, dbUrl, ctx)

	if err != nil {
		return nil, err
	}

	res, err := data.SyncGithubInstallationDataToTenant(
		installationId,
		installationOrgName,
		installationOrgId,
		authId,
		db, ctx,
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
	client, err := activity.NewInstallationClient(installationId)

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
	client, err := activity.NewInstallationClient(installationId)

	if err != nil {
		fmt.Printf("Could not create new installation client. Error: %v", err.Error())
		return nil, err
	}
	return data.GetInstallationTeamMembers(ctx, installationOrgName, teamSlug, client)
}

func (activity *GithubActivities) GetInstallationTeamMembersWithEmails(ctx context.Context, installationId int64, members []*github.User) (data.ExtendedMembers, error) {

	client, err := activity.NewInstallationClient(installationId)

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
	db, err := activity.GetCachedTenantDB(&activity.Connections, dbUrl, ctx)

	if err != nil {
		return false, err
	}

	return data.SyncTeamsAndMembersToTenant(
		ctx,
		teamWithMembers,
		dbUrl,
		githubOrganizationId,
		organizationId,
		db)
}

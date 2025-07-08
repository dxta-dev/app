package activity

import (
	"context"
	"errors"

	"github.com/dxta-dev/app/internal/onboarding"
)

type GithubInstallationActivities struct {
	githubAppClient onboarding.GithubAppClient
}

func NewGithubInstallationActivities(GithubAppClient onboarding.GithubAppClient) *GithubInstallationActivities {
	return &GithubInstallationActivities{
		githubAppClient: GithubAppClient,
	}
}

type GithubInstallationOrganization struct {
	OrganizationID    int64
	OrganizationLogin string
}

func (gia *GithubInstallationActivities) GetInstallationOrganization(
	ctx context.Context,
	installationId int64,
) (*GithubInstallationOrganization, error) {
	account, err := gia.githubAppClient.GetInstallationAccount(ctx, installationId)

	if err != nil {
		return nil, errors.New("failed to retrieve installation account: " + err.Error())
	}

	return &GithubInstallationOrganization{
		OrganizationID:    account.ID,
		OrganizationLogin: account.Login,
	}, nil
}

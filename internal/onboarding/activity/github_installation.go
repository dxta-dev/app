package activity

import (
	"context"
	"fmt"

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

func (gia *GithubInstallationActivities) GetGithubInstallation(
	ctx context.Context,
	installationId int64,
) (string, error) {
	login, err := gia.githubAppClient.GetOrganizationLogin(ctx, installationId)

	if err != nil {
		fmt.Printf("Could not retrieve installation. Error: %v", err.Error())
		return "", err
	}

	return login, nil
}

package activity

import (
	"context"
	"fmt"

	"github.com/dxta-dev/app/internal/onboarding"
)

type GithubActivities struct {
	githubAppClient onboarding.GithubAppClient
}

func NewGithubInstallationActivities(GithubAppClient onboarding.GithubAppClient) *GithubActivities {
	return &GithubActivities{
		githubAppClient: GithubAppClient,
	}
}

func (ga *GithubActivities) GetGithubInstallation(
	ctx context.Context,
	installationId int64,
) (string, error) {
	login, err := ga.githubAppClient.GetOrganizationLogin(ctx, installationId)

	if err != nil {
		fmt.Printf("Could not retrieve installation. Error: %v", err.Error())
		return "", err
	}

	return login, nil
}

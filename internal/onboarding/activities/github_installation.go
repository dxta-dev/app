package activity

import (
	"context"
	"fmt"

	"github.com/dxta-dev/app/internal/onboarding"
	"github.com/google/go-github/v73/github"
)

type InstallationActivities struct {
	GithubAppClient onboarding.AppClient
}

func GithubInstallationActivities(GithubAppClient onboarding.AppClient) *InstallationActivities {
	return &InstallationActivities{
		GithubAppClient: GithubAppClient,
	}
}

func (a InstallationActivities) GetGithubInstallation(
	ctx context.Context,
	installationId int64,
) (*github.Installation, error) {
	installation, _, err := a.GithubAppClient.GetInstallation(ctx, installationId)

	if err != nil {
		fmt.Printf("Could not retrieve installation. Error: %v", err.Error())
		return nil, err
	}

	return installation, nil
}

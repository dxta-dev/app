package activity

import (
	"context"
	"fmt"

	"github.com/google/go-github/v73/github"
)

type InstallationActivities struct {
	GithubAppClient *github.Client
}

func GithubInstallationActivities(GithubAppClient *github.Client) *InstallationActivities {
	return &InstallationActivities{
		GithubAppClient: GithubAppClient,
	}
}

func (a InstallationActivities) GetGithubInstallation(
	ctx context.Context,
	installationId int64,
) (*github.Installation, error) {
	installation, _, err := a.GithubAppClient.Apps.GetInstallation(ctx, installationId)

	if err != nil {
		fmt.Printf("Could not retrieve installation. Error: %v", err.Error())
		return nil, err
	}

	return installation, nil
}

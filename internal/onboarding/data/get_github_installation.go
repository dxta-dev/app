package data

import (
	"context"
	"fmt"

	"github.com/google/go-github/v72/github"
)

func GetGithubInstallation(installationId int64, githubAppClient *github.Client, ctx context.Context) (*github.Installation, error) {
	installation, _, err := githubAppClient.Apps.GetInstallation(ctx, installationId)

	if err != nil {
		fmt.Printf("Could not retrieve installation. Error: %v", err.Error())
		return nil, err
	}

	return installation, nil
}

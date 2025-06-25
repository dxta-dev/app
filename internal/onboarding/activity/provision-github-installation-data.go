package activity

import (
	"context"

	"github.com/dxta-dev/app/internal/onboarding/data"
	"github.com/google/go-github/v72/github"
)

func GetGithubInstallation(ctx context.Context, installationId int64) (*github.Installation, error) {
	installations, err := data.GithubConfig.GetGithubInstallation(installationId, ctx)

	if err != nil {
		return nil, err
	}

	return installations, nil
}

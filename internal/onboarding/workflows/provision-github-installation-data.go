package workflows

import (
	"context"
	"fmt"
	"time"

	api "github.com/dxta-dev/app/internal/internal-api"
	"github.com/dxta-dev/app/internal/onboarding/activity"
	"github.com/google/go-github/v72/github"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

func ProvisionGithubInstallationData(ctx workflow.Context, installationId int64) (*github.Installation, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// 1. Get installation
	var installation *github.Installation
	err := workflow.ExecuteActivity(ctx, activity.GetGithubInstallation, installationId).Get(ctx, &installation)

	if err != nil {
		return nil, err
	}

	// TO - DO
	// 2. Add installation data to tenant
	// 3. Add Teams to tenant

	return installation, nil
}

type Args struct {
	TemporalOnboardingQueueName string
	InstallationId              int64
	ApiState                    api.State
}

func ExecuteGithubInstallationDataProvision(
	ctx context.Context,
	temporalClient client.Client,
	args Args,
) (*github.Installation, error) {

	wr, err := temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:        fmt.Sprintf("onboarding-workflow-%v", time.Now().Format("20060102150405")),
			TaskQueue: args.TemporalOnboardingQueueName,
		},
		ProvisionGithubInstallationData,
		args.InstallationId,
	)
	if err != nil {
		return nil, err
	}

	var installation *github.Installation

	err = wr.Get(ctx, &installation)

	if err != nil {
		return nil, err
	}

	return installation, nil
}

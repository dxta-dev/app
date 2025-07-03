package workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v73/github"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/dxta-dev/app/internal/onboarding/activity"
)

type AfterGithubInstallationParams struct {
	InstallationID int64
	AuthID         string
	DBURL          string
}

func AfterGithubInstallationWorkflow(
	ctx workflow.Context,
	params AfterGithubInstallationParams,
) (err error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 30,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 10,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	installationId := params.InstallationID

	var installation *github.Installation
	err = workflow.ExecuteActivity(ctx, (*activity.GithubInstallationActivities).GetGithubInstallation, installationId).
		Get(ctx, &installation)

	if err != nil {
		return
	}

	return
}

type ExecuteAfterGithubInstallationParams struct {
	TemporalOnboardingQueueName string
	InstallationID              int64
	AuthID                      string
	DBURL                       string
}

func ExecuteAfterGithubInstallationWorkflow(
	ctx context.Context,
	temporalClient client.Client,
	params ExecuteAfterGithubInstallationParams,
) (string, error) {
	_, err := temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID: fmt.Sprintf(
				"onboarding-workflow-github-%v-%v",
				params.InstallationID, params.AuthID,
			),
			TaskQueue: params.TemporalOnboardingQueueName,
		},
		AfterGithubInstallationWorkflow,
		AfterGithubInstallationParams{
			InstallationID: params.InstallationID,
			AuthID:         params.AuthID,
			DBURL:          params.DBURL,
		},
	)

	if err != nil {
		return "Unable to execute ", err
	}

	return "Success", nil
}

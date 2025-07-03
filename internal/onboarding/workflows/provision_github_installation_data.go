package workflows

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v73/github"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	activity "github.com/dxta-dev/app/internal/onboarding/activities"
)

type ProvisionGithubInstallationDataParams struct {
	InstallationId int64
	AuthId         string
	DBUrl          string
}

func ProvisionGithubInstallationData(ctx workflow.Context, params ProvisionGithubInstallationDataParams) (count int, err error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 30,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 10,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	installationId := params.InstallationId
	/* authId := params.AuthId
	dbUrl := params.DBUrl */

	var installation *github.Installation
	err = workflow.ExecuteActivity(ctx, (*activity.InstallationActivities).GetGithubInstallation, installationId).Get(ctx, &installation)

	if err != nil {
		return
	}

	return
}

type ExecuteGithubInstallationDataProvisionParams struct {
	TemporalOnboardingQueueName string
	InstallationId              int64
	AuthId                      string
	DBUrl                       string
}

func ExecuteGithubInstallationDataProvision(
	ctx context.Context,
	temporalClient client.Client,
	params ExecuteGithubInstallationDataProvisionParams,
) (string, error) {
	_, err := temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:        fmt.Sprintf("onboarding-workflow-%v", time.Now().Format("20060102150405")),
			TaskQueue: params.TemporalOnboardingQueueName,
		},
		ProvisionGithubInstallationData,
		ProvisionGithubInstallationDataParams{
			InstallationId: params.InstallationId,
			AuthId:         params.AuthId,
			DBUrl:          params.DBUrl,
		},
	)

	if err != nil {
		return "Unable to execute ", err
	}

	return "Success", nil
}

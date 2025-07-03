package workflows

import (
	"context"
	"fmt"
	"time"

	activity "github.com/dxta-dev/app/internal/onboarding/activities"
	"go.temporal.io/sdk/workflow"

	"github.com/dxta-dev/app/internal/onboarding"
	"go.temporal.io/sdk/client"
)

func CountUsers(ctx workflow.Context) (int, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	var count int
	err := workflow.ExecuteActivity(ctx, (*activity.UserActivites).CountUsers).Get(ctx, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func ExecuteCountUsersWorkflow(
	ctx context.Context,
	temporalClient client.Client,
	cfg onboarding.Config,
) (int, error) {
	wr, err := temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:        fmt.Sprintf("count-users-%v", time.Now().Format("20060102150405")),
			TaskQueue: cfg.TemporalOnboardingQueueName,
		},
		CountUsers,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to start CountUsersWorkflow: %w", err)
	}

	var result int
	err = wr.Get(ctx, &result)
	if err != nil {
		return 0, fmt.Errorf("failed to get result from CountUsersWorkflow: %w", err)
	}

	return result, nil
}

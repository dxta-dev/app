package workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/dxta-dev/app/internal/onboarding/activity"
	"go.temporal.io/sdk/workflow"

	"github.com/dxta-dev/app/internal/onboarding"
	"go.temporal.io/sdk/client"
)

func CountUsersWorkflow(ctx workflow.Context, dsn string) (int, error) {
	ao := workflow.ActivityOptions{}

	ctx = workflow.WithActivityOptions(ctx, ao)

	var count int
	err := workflow.ExecuteActivity(ctx, activity.CountUsersActivity, dsn).Get(ctx, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func ExecuteCountUsersWorkflow(ctx context.Context, temporalClient client.Client, cfg onboarding.Config) (int, error) {
	wr, err := temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:        fmt.Sprintf("count-users-workflow-%v", time.Now().Format("20060102150405")),
		TaskQueue: cfg.TemporalQueueName,
	}, CountUsersWorkflow, cfg.UsersDSN)
	if err != nil {
		return 0, err
	}

	var result int
	wr.Get(ctx, &result)

	return result, nil
}

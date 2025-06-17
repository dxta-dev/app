package workflow

import (
	"log"

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

func ExecuteCountUsersWorkflow(temporalClient client.Client, cfg onboarding.Config) {
	log.Fatal("Not implemented yet")
}

package workflow

import (
	"github.com/dxta-dev/app/internal/onboarding/activity"
	"go.temporal.io/sdk/workflow"
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

package workflow

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dxta-dev/app/internal/onboarding/activity"
	"github.com/dxta-dev/app/internal/util"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type CreateTenantDBParams struct {
	DBName           string
	AuthID           string
	OrganizationName string
}

func CreateTenantDBWorkflow(
	ctx workflow.Context,
	params CreateTenantDBParams,
) (err error) {
	if params.DBName == "" || params.AuthID == "" || params.OrganizationName == "" {
		err = errors.New("bad request")
		return
	}

	sanitizedDBName := util.SanitizeString(params.DBName)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 30,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 10,
			InitialInterval: time.Millisecond * 500,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	var newDBData activity.CreateTenantDBRes

	err = workflow.ExecuteActivity(
		ctx,
		(*activity.CreateTenantActivities).CreateTenantDB,
		sanitizedDBName,
	).Get(ctx, &newDBData)

	if err != nil {
		return
	}

	mapTenantFuture := workflow.ExecuteActivity(
		ctx,
		(*activity.CreateTenantActivities).AddTenantDBToMap,
		params.AuthID,
		params.DBName,
		newDBData.Database.DBURL,
		newDBData.Database.Name,
	)

	upsertTenantInfoFuture := workflow.ExecuteActivity(
		ctx,
		(*activity.TenantActivities).UpsertTenantDBInfo,
		params.DBName,
		newDBData.Database.DBURL,
		newDBData.Database.Name,
	)

	createOrganizationFuture := workflow.ExecuteActivity(
		ctx,
		(*activity.TenantActivities).CreateOrganization,
		params.OrganizationName,
		params.AuthID,
		newDBData.Database.DBURL,
	)

	var tenantDBMapRes bool

	err = mapTenantFuture.Get(ctx, &tenantDBMapRes)

	if err != nil {
		return
	}

	var upsertTenantInfoRes bool

	err = upsertTenantInfoFuture.Get(ctx, &upsertTenantInfoRes)

	if err != nil {
		return
	}

	var createOrganizationRes bool

	err = createOrganizationFuture.Get(ctx, &createOrganizationRes)

	if err != nil {
		return
	}

	return
}

type ExecuteCreateTenantDBWorkflowParams struct {
	TemporalOnboardingQueueName string
	AuthID                      string
	DBName                      string
	OrganizationName            string
}

func ExecuteCreateTenantDBWorkflow(
	ctx context.Context,
	temporalClient client.Client,
	params ExecuteCreateTenantDBWorkflowParams,
) (string, error) {
	_, err := temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID: fmt.Sprintf(
				"onboarding-workflow-github-%v-%v",
				params.AuthID,
				params.DBName,
			),
			TaskQueue: params.TemporalOnboardingQueueName,
		},
		CreateTenantDBWorkflow,
		CreateTenantDBParams{
			AuthID:           params.AuthID,
			DBName:           params.DBName,
			OrganizationName: params.OrganizationName,
		},
	)

	if err != nil {
		return "Unable to execute ", err
	}

	return "Success", nil
}

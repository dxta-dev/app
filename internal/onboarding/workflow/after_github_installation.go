package workflow

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/dxta-dev/app/internal/onboarding"
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
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 10,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	var installation activity.GithubInstallationOrganization

	err = workflow.ExecuteActivity(
		ctx,
		(*activity.GithubInstallationActivities).GetInstallationOrganization,
		params.InstallationID,
	).Get(ctx, &installation)

	if err != nil {
		return
	}

	var organizationId int64

	err = workflow.ExecuteActivity(
		ctx,
		(*activity.TenantActivities).GetOrganizationIDByAuthID,
		params.AuthID,
		params.DBURL,
	).Get(ctx, &organizationId)

	if err != nil {
		return
	}

	var githubOrganizationId int64

	err = workflow.ExecuteActivity(
		ctx,
		(*activity.TenantActivities).UpsertGithubOrganization,
		params.DBURL,
		params.InstallationID,
		installation.OrganizationLogin,
		installation.OrganizationID,
		organizationId,
	).Get(ctx, &githubOrganizationId)

	if err != nil {
		return
	}

	var githubTeams onboarding.Teams

	err = workflow.ExecuteActivity(
		ctx,
		(*activity.GithubActivities).GetTeams,
		params.InstallationID,
		installation.OrganizationLogin,
	).Get(ctx, &githubTeams)

	if err != nil {
		return
	}

	counter := 0

	for _, team := range githubTeams {
		workflow.Go(ctx, func(gctx workflow.Context) {

			if err != nil {
				return
			}

			var teamMembers onboarding.Members

			err = workflow.ExecuteActivity(
				gctx,
				(*activity.GithubActivities).GetTeamMembers,
				params.InstallationID,
				installation.OrganizationLogin,
				team.Slug,
			).Get(gctx, &teamMembers)

			if err != nil {
				return
			}

			var teamMembersFutures []workflow.Future

			for _, member := range teamMembers {
				teamMembersFutures = append(teamMembersFutures, workflow.ExecuteActivity(gctx, (*activity.GithubActivities).GetTeamMemberEmail, params.InstallationID, member))
			}

			var teamMembersWithEmails onboarding.Members

			for i := 0; i < len(teamMembers); i++ {
				var memberWithEmail onboarding.Member

				err := teamMembersFutures[i].Get(gctx, &memberWithEmail)

				if err != nil {
					return
				}

				teamMembersWithEmails = append(teamMembersWithEmails, memberWithEmail)
			}

			var teamAndMembersUpsertRes *bool

			err = workflow.ExecuteActivity(
				gctx,
				(*activity.TenantActivities).UpsertTeamAndMembers,
				team,
				teamMembersWithEmails,
				params.DBURL,
				githubOrganizationId,
				organizationId,
			).Get(gctx, &teamAndMembersUpsertRes)

			if err != nil {
				return
			}

			// Count number of finished go routines
			// so we can unblock calling thread when
			// all go routines finish
			counter += 1
		})
	}

	_ = workflow.Await(ctx, func() bool {
		return err != nil || counter == len(githubTeams)
	})

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

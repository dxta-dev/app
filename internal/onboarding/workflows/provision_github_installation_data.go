package workflows

import (
	"context"
	"fmt"
	"time"

	"github.com/dxta-dev/app/internal/onboarding/activities"
	"github.com/dxta-dev/app/internal/onboarding/data"
	"github.com/google/go-github/v72/github"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type GithubDataProvisionResponse struct {
	Installation *github.Installation   `json:"installation"`
	Teams        []data.TeamWithMembers `json:"teams"`
}

func ProvisionGithubInstallationData(ctx workflow.Context, installationId int64, authId string, dbUrl string) (count int, err error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// 1. Get installation data
	var installation *github.Installation
	err = workflow.ExecuteActivity(ctx, (*activities.GithubActivities).GetGithubInstallation, installationId).Get(ctx, &installation)

	if err != nil {
		return
	}

	// 2. Store installation data to tenant
	var syncResult *data.SyncGithubDataResult
	err = workflow.ExecuteActivity(ctx, (*activities.DBActivities).SyncGithubInstallationDataToTenant, installationId, installation.Account.Login, installation.Account.ID, authId, dbUrl).Get(ctx, &syncResult)

	if err != nil {
		return
	}

	// 3. Retrieve installation Github Teams and Github Members and store them in tenant db and create copy of them as DXTA Teams and Members
	if installation.TargetType != nil && *installation.TargetType == "Organization" {
		// 3.1 Retrieve all installation Github Teams
		var teams []*github.Team
		err = workflow.ExecuteActivity(ctx, (*activities.GithubActivities).GetInstallationTeams, installation.Account.Login, installationId).Get(ctx, &teams)

		if err != nil {
			return
		}

		// 3.2 Retrieve Github Members and store installation Github Teams and Github Members to tenant db
		for _, team := range teams {
			workflow.Go(ctx, func(gctx workflow.Context) {
				teamWithMembers := data.TeamWithMembers{Team: team, Members: data.ExtendedMembers{}}

				var members []*github.User

				err = workflow.ExecuteActivity(gctx, (*activities.GithubActivities).GetInstallationTeamMembers, installationId, installation.Account.Login, team.Slug).Get(gctx, &members)

				if err != nil {
					return
				}

				var membersWithEmails *data.ExtendedMembers

				err = workflow.ExecuteActivity(gctx, (*activities.GithubActivities).GetInstallationTeamMembersWithEmails, installationId, members).Get(gctx, &membersWithEmails)

				if err != nil {
					return
				}

				teamWithMembers.Members = *membersWithEmails

				var syncTeamsAndMembersRes *bool
				err = workflow.ExecuteActivity(gctx, (*activities.DBActivities).SyncTeamsAndMembersToTenant, teamWithMembers, dbUrl, syncResult.GithubOrganizationId, syncResult.OrganizationId).Get(gctx, &syncTeamsAndMembersRes)

				if err != nil {
					return
				}

				// Count number of finished go routines
				// so we can unblock calling thread when
				// all go routines finish
				count += 1
			})
		}

		_ = workflow.Await(ctx, func() bool {
			return err != nil || count == len(teams)
		})
	}

	return
}

type Args struct {
	TemporalOnboardingQueueName string
	InstallationId              int64
	AuthId                      string
	DBUrl                       string
}

func ExecuteGithubInstallationDataProvision(
	ctx context.Context,
	temporalClient client.Client,
	args Args,
) (string, error) {
	_, err := temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:        fmt.Sprintf("onboarding-workflow-%v", time.Now().Format("20060102150405")),
			TaskQueue: args.TemporalOnboardingQueueName,
		},
		ProvisionGithubInstallationData,
		args.InstallationId,
		args.AuthId,
		args.DBUrl,
	)

	if err != nil {
		return "Unable to execute ", err
	}

	return "Success", nil
}

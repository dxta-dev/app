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

func ProvisionGithubInstallationData(ctx workflow.Context, installationId int64, authId string, dbUrl string) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// 1. Get installation data
	var installation *github.Installation
	err := workflow.ExecuteActivity(ctx, (*activities.GithubActivities).GetGithubInstallation, installationId).Get(ctx, &installation)

	if err != nil {
		return err
	}

	// 2. Store installation data to tenant
	var syncResult *data.SyncGithubDataResult
	err = workflow.ExecuteActivity(ctx, (*activities.DBActivities).SyncGithubInstallationDataToTenant, installationId, installation.Account.Login, installation.Account.ID, authId, dbUrl).Get(ctx, &syncResult)

	if err != nil {
		return err
	}

	// 3. Add Teams to tenant
	if installation.TargetType != nil && *installation.TargetType == "Organization" {
		// 3.1 Retrieve all teams
		var teams []*github.Team
		err := workflow.ExecuteActivity(ctx, (*activities.GithubActivities).GetInstallationTeams, installation.Account.Login, installationId).Get(ctx, &teams)

		if err != nil {
			return err
		}

		// 3.2 Retrieve members and store teams and members to tenant db
		for _, t := range teams {
			team := t
			teamWithMembers := data.TeamWithMembers{Team: team, Members: data.ExtendedMembers{}}

			var members []*github.User
			err := workflow.ExecuteActivity(ctx, (*activities.GithubActivities).GetInstallationTeamMembers, installationId, installation.Account.Login, team.Slug).Get(ctx, &members)

			if err != nil {
				return err
			}

			var membersWithEmails *data.ExtendedMembers

			err = workflow.ExecuteActivity(ctx, (*activities.GithubActivities).GetInstallationTeamMembersWithEmails, installationId, members).Get(ctx, &membersWithEmails)

			if err != nil {
				return err
			}

			teamWithMembers.Members = *membersWithEmails

			var syncTeamsAndMembersRes *bool
			err = workflow.ExecuteActivity(ctx, (*activities.DBActivities).SyncTeamsAndMembersToTenant, teamWithMembers, dbUrl, syncResult.GithubOrganizationId, syncResult.OrganizationId).Get(ctx, &syncTeamsAndMembersRes)

			if err != nil {
				return err
			}
		}
	}

	return nil
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

	return "success", nil
}

package workflows

import (
	"context"
	"fmt"
	"time"

	"github.com/dxta-dev/app/internal/onboarding/activities"
	"github.com/google/go-github/v72/github"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type TeamWithMembers struct {
	Team    *github.Team
	Members []*github.User
}

type GithubDataProvisionResponse struct {
	Installation *github.Installation `json:"installation"`
	Teams        []TeamWithMembers    `json:"teams"`
}

func ProvisionGithubInstallationData(ctx workflow.Context, installationId int64, authId string, dbUrl string) (*GithubDataProvisionResponse, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// 1. Get installation
	var installation *github.Installation
	err := workflow.ExecuteActivity(ctx, (*activities.GithubActivities).GetGithubInstallation, installationId).Get(ctx, &installation)

	if err != nil {
		return nil, err
	}

	// 2. Add installation data to tenant
	var syncResult bool
	err = workflow.ExecuteActivity(ctx, (*activities.DBActivities).SyncGithubInstallationDataToTenant, installationId, installation.Account.Login, installation.Account.ID, authId, dbUrl).Get(ctx, &syncResult)

	if err != nil {
		return nil, err
	}

	// 3. Add Teams to tenant
	var teamsWithMembers []TeamWithMembers
	if installation.TargetType != nil && *installation.TargetType == "Organization" {
		var teams []*github.Team
		// 3.1 Retrieve all teams and its members
		err := workflow.ExecuteActivity(ctx, (*activities.GithubActivities).GetInstallationTeams, installation.Account.Login, installationId).Get(ctx, &teams)

		if err != nil {
			return nil, err
		}

		// 3.2 Retrieve members for each team
		// TO-DO Run this in go routines instead in sequence
		for _, team := range teams {
			teamWithMembers := TeamWithMembers{Team: team, Members: []*github.User{}}

			var members []*github.User

			err := workflow.ExecuteActivity(ctx, (*activities.GithubActivities).GetInstallationTeamMembers, installationId, installation.Account.Login, *team.Slug).Get(ctx, &members)

			if err != nil {
				return nil, err
			}

			teamWithMembers.Members = members
			teamsWithMembers = append(teamsWithMembers, teamWithMembers)
		}

	}

	return &GithubDataProvisionResponse{Installation: installation, Teams: teamsWithMembers}, nil
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
) (*GithubDataProvisionResponse, error) {
	wr, err := temporalClient.ExecuteWorkflow(
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
		return nil, err
	}

	var installation *GithubDataProvisionResponse

	err = wr.Get(ctx, &installation)

	if err != nil {
		return nil, err
	}

	return installation, nil
}

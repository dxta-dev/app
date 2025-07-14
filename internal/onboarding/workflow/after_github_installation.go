package workflow

import (
	"context"
	"errors"
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

func addMemberToMap(
	team onboarding.Team,
	member onboarding.ExtendedMember,
	membersMap activity.MembersRecordMap,
) activity.MemberRecord {
	m, ok := membersMap[*member.Login]

	if ok {
		m.Teams = append(m.Teams, struct {
			Name   *string
			TeamID *int64
		}{Name: team.Name})

		membersMap[*member.Login] = m

	} else {
		t := append(m.Teams, struct {
			Name   *string
			TeamID *int64
		}{Name: team.Name})

		membersMap[*member.Login] = activity.MemberRecord{
			ID:    member.ID,
			Login: member.Login,
			Email: member.Email,
			Name:  member.Name,
			Teams: t,
		}
	}
	return m
}

func AfterGithubInstallationWorkflow(
	ctx workflow.Context,
	params AfterGithubInstallationParams,
) (err error) {

	if params.InstallationID == 0 || params.AuthID == "" || params.DBURL == "" {
		err = errors.New("bad request")
		return
	}

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

	teamsMap := activity.TeamsRecordMap{}
	membersMap := activity.MembersRecordMap{}

	for _, team := range githubTeams {
		workflow.Go(ctx, func(gctx workflow.Context) {
			var teamMembers onboarding.Members

			err = workflow.ExecuteActivity(
				gctx,
				(*activity.GithubActivities).GetTeamMembers,
				params.InstallationID,
				installation.OrganizationLogin,
				team.Slug,
				team.Name,
			).Get(gctx, &teamMembers)

			if err != nil {
				return
			}

			var teamMembersFutures []workflow.Future

			for _, member := range teamMembers {
				teamMembersFutures = append(
					teamMembersFutures,
					workflow.ExecuteActivity(
						gctx,
						(*activity.GithubActivities).GetExtendedTeamMember,
						params.InstallationID,
						member,
					))
			}

			for i := range teamMembers {
				var member onboarding.ExtendedMember

				err := teamMembersFutures[i].Get(gctx, &member)

				if err != nil {
					return
				}

				addMemberToMap(team, member, membersMap)
			}

			teamsMap[*team.Name] = activity.TeamsRecord{
				ID:           team.ID,
				Name:         team.Name,
				GithubTeamID: nil,
				TeamID:       nil,
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

	err = workflow.ExecuteActivity(
		ctx,
		(*activity.TenantActivities).UpsertTeams,
		params.DBURL,
		githubOrganizationId,
		organizationId,
		teamsMap,
	).Get(ctx, &teamsMap)

	if err != nil {
		return
	}

	var newGithubMembers activity.MembersRecordMap

	err = workflow.ExecuteActivity(
		ctx,
		(*activity.TenantActivities).UpsertGithubMembers,
		params.DBURL,
		membersMap,
		teamsMap,
	).Get(ctx, &newGithubMembers)

	if err != nil {
		return
	}

	if len(newGithubMembers) > 0 {
		newMembers := make([]activity.MemberRecord, 0)

		err = workflow.ExecuteActivity(
			ctx,
			(*activity.TenantActivities).CreateTeamMembers,
			params.DBURL,
			newGithubMembers,
			organizationId,
		).Get(ctx, &newMembers)

		var joinRes bool

		err = workflow.ExecuteActivity(
			ctx,
			(*activity.TenantActivities).JoinTeamsMembers,
			params.DBURL,
			newMembers,
		).Get(ctx, &joinRes)

		if err != nil {
			return
		}
	}

	return
}

type ExecuteAfterGithubInstallationParams struct {
	TemporalOnboardingQueueName string
	InstallationID              int64
	AuthID                      string
	DBURL                       string
	DBDomainName                string
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
				params.DBDomainName,
				params.InstallationID,
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

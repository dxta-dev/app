package onboarding

import (
	"context"
	"errors"

	"github.com/google/go-github/v73/github"
)

type Account struct {
	ID    int64
	Login string
}

func (gac *GithubAppClient) GetInstallationAccount(
	ctx context.Context,
	installationID int64,
) (*Account, error) {
	installation, _, error := gac.client.Apps.GetInstallation(ctx, installationID)
	if error != nil {
		return nil, errors.New("failed to get installation: " + error.Error())

	}

	if installation.Account == nil || installation.Account.Login == nil {
		return nil, errors.New("installation account or login is nil")
	}

	if installation.TargetType == nil || *installation.TargetType != "Organization" {
		return nil, errors.New("installation is not for an organization")
	}

	return &Account{
		ID:    *installation.Account.ID,
		Login: *installation.Account.Login,
	}, nil
}

type ListOptions struct {
	PerPage int
	//Limit to max amount of pages we want to visit
	MaxPages int
}

type Team struct {
	Name *string `json:"name,omitempty"`
	ID   *int64  `json:"id,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

type Teams []Team

func (gic *GithubInstallationClient) GetTeams(
	ctx context.Context,
	organizationName string,
	listOptions *ListOptions,
) (Teams, error) {
	opts := &github.ListOptions{}

	if listOptions != nil {
		opts = &github.ListOptions{PerPage: listOptions.PerPage}
	}

	var allTeams Teams

	for range listOptions.MaxPages {
		teams, res, err := gic.client.Teams.ListTeams(ctx, organizationName, opts)

		if err != nil {
			return nil, errors.New("failed to retrieve github teams list: " + err.Error())
		}

		t := make(Teams, 0)

		for _, team := range teams {
			t = append(t, Team{Name: team.Name, ID: team.ID, Slug: team.Slug})
		}

		allTeams = append(allTeams, t...)

		if res.NextPage == 0 {
			break
		}

		opts.Page = res.NextPage
	}

	return allTeams, nil
}

type Member struct {
	Login *string `json:"login,omitempty"`
	ID    *int64  `json:"id,omitempty"`
}

type Members []Member

func (gic *GithubInstallationClient) GetTeamMembers(
	ctx context.Context,
	organizationName string,
	teamSlug string,
	teamName string,
	listOptions *ListOptions,
) (Members, error) {
	lo := github.ListOptions{}

	if listOptions != nil {
		lo = github.ListOptions{PerPage: listOptions.PerPage}
	}

	opts := &github.TeamListTeamMembersOptions{ListOptions: lo}

	var allMembers Members

	for range listOptions.MaxPages {
		members, res, err := gic.client.Teams.ListTeamMembersBySlug(
			ctx,
			organizationName,
			teamSlug,
			opts,
		)

		if err != nil {
			return nil, errors.New("failed to retrieve github team members list: " + err.Error())
		}

		m := make(Members, 0)

		for _, member := range members {
			m = append(m,
				Member{
					Login: member.Login,
					ID:    member.ID})
		}

		allMembers = append(allMembers, m...)

		if res.NextPage == 0 {
			break
		}

		opts.Page = res.NextPage
	}

	return allMembers, nil
}

type ExtendedMember struct {
	*Member
	Email *string
	Name  *string
}

func (gic *GithubInstallationClient) GetExtendedTeamMember(
	ctx context.Context,
	teamMember Member,
) (*ExtendedMember, error) {
	user, _, err := gic.client.Users.Get(ctx, *teamMember.Login)

	if err != nil {
		return nil, errors.New("failed to retrieve github team member email: " + err.Error())
	}

	return &ExtendedMember{
		Member: &teamMember,
		Email:  user.Email,
		Name:   user.Name,
	}, nil
}

package data

import (
	"context"
	"fmt"

	"github.com/google/go-github/v72/github"
)

func GetInstallationTeamMembers(ctx context.Context, installationOrgName string, teamSlug string, client *github.Client, extendWithEmail bool) ([]*github.User, error) {
	opts := &github.TeamListTeamMembersOptions{ListOptions: github.ListOptions{PerPage: 100}}

	var allMembers []*github.User

	for {
		members, res, err := client.Teams.ListTeamMembersBySlug(ctx, installationOrgName, teamSlug, opts)

		if err != nil {
			fmt.Printf("Could not retrieve installation. Error: %v", err.Error())
			return nil, err
		}

		// TO-DO Handle if we need to extend member with email.
		// For each member request towards github is needed so it
		// makes sense to run each request in its own go routine

		allMembers = append(allMembers, members...)

		if res.NextPage == 0 {
			break
		}

		opts.Page = res.NextPage
	}
	return allMembers, nil
}

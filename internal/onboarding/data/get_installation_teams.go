package data

import (
	"context"
	"fmt"

	"github.com/google/go-github/v72/github"
)

type TeamWithMembers struct {
	Team    *github.Team
	Members ExtendedMembers
}

func GetInstallationTeams(
	ctx context.Context,
	installationOrgName string,
	client *github.Client,
) ([]*github.Team, error) {

	opt := &github.ListOptions{PerPage: 100}

	var allTeams []*github.Team

	for {
		teams, res, err := client.Teams.ListTeams(ctx, installationOrgName, opt)

		if err != nil {
			fmt.Printf("Could not retrieve installation. Error: %v", err.Error())
			return nil, err
		}

		allTeams = append(allTeams, teams...)

		if res.NextPage == 0 {
			break
		}

		opt.Page = res.NextPage
	}

	return allTeams, nil

}

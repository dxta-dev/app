package data

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/go-github/v72/github"
	"golang.org/x/sync/errgroup"
)

type ExtendedMember struct {
	*github.User
	Email *string `json:"email,omitempty"`
	Name  *string `json:"name,omitempty"`
}

type ExtendedMembers []ExtendedMember

func GetInstallationTeamMembers(ctx context.Context, installationOrgName string, teamSlug string, client *github.Client) ([]*github.User, error) {
	opts := &github.TeamListTeamMembersOptions{ListOptions: github.ListOptions{PerPage: 100}}

	var allMembers []*github.User

	for {
		members, res, err := client.Teams.ListTeamMembersBySlug(ctx, installationOrgName, teamSlug, opts)

		if err != nil {
			fmt.Printf("Could not retrieve installation. Error: %v", err.Error())
			return nil, err
		}

		allMembers = append(allMembers, members...)

		if res.NextPage == 0 {
			break
		}

		opts.Page = res.NextPage
	}

	return allMembers, nil
}

type AllMembersContainer struct {
	mu         sync.Mutex
	allMembers ExtendedMembers
}

func (amc *AllMembersContainer) extendMember(member *github.User, Email *string, Name *string) {
	amc.mu.Lock()
	defer amc.mu.Unlock()
	amc.allMembers = append(amc.allMembers, ExtendedMember{User: member, Email: Email, Name: Name})
}

func GetInstallationTeamMembersWithEmails(ctx context.Context, members []*github.User, client *github.Client) (ExtendedMembers, error) {
	c := AllMembersContainer{
		allMembers: ExtendedMembers{},
	}

	g := new(errgroup.Group)

	for _, m := range members {

		g.Go(func() error {
			user, _, err := client.Users.Get(ctx, *m.Login)

			if err != nil {
				return err
			}

			c.extendMember(m, user.Email, user.Name)
			return nil
		})

	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return c.allMembers, nil
}

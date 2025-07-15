package activity

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dxta-dev/app/internal/onboarding"
)

type GithubActivities struct {
	GithubConfig onboarding.GithubConfig
}

func NewGithubActivities(githubConfig onboarding.GithubConfig) *GithubActivities {
	return &GithubActivities{
		GithubConfig: githubConfig,
	}
}

func (ga *GithubActivities) GetTeams(
	ctx context.Context,
	installationId int64,
	organizationName string,
) (onboarding.Teams, error) {

	client, err := onboarding.GetCachedGithubInstallationClient(installationId, ga.GithubConfig)

	if err != nil {
		return nil, errors.New("could not create new installation client to get github teams " + err.Error())
	}

	return client.GetTeams(ctx, organizationName, &onboarding.ListOptions{PerPage: 100, MaxPages: 50})
}

func (ga *GithubActivities) GetTeamMembers(
	ctx context.Context,
	installationId int64,
	organizationName string,
	teamSlug string,
	teamName string,
) (onboarding.Members, error) {
	client, err := onboarding.GetCachedGithubInstallationClient(installationId, ga.GithubConfig)

	if err != nil {
		return nil, errors.New("could not create new installation client to get github members " + err.Error())
	}

	return client.GetTeamMembers(
		ctx,
		organizationName,
		teamSlug,
		teamName,
		&onboarding.ListOptions{PerPage: 100, MaxPages: 50})
}

func (ga *GithubActivities) GetExtendedTeamMember(
	ctx context.Context,
	installationId int64,
	teamMember onboarding.Member,
	teamSlug string,
) (*onboarding.ExtendedMember, error) {
	client, err := onboarding.GetCachedGithubInstallationClient(installationId, ga.GithubConfig)

	if err != nil {
		return nil, errors.New("could not create new installation client to get github members " + err.Error())
	}

	return client.GetExtendedTeamMember(ctx, teamMember)
}

type TeamsRecord struct {
	ID   *int64
	Name *string
	// ID of a record after insertion to github_teams table
	GithubTeamID *int64
	// ID of a record after insertion to teams table
	TeamID *int64
}

type TeamsRecordMap map[string]TeamsRecord

func (ta *TenantActivities) UpsertTeams(
	ctx context.Context,
	DBURL string,
	githubOrganizationId int64,
	organizationId int64,
	teamsRecordMap *TeamsRecordMap,
) (res *TeamsRecordMap, err error) {
	db, err := ta.GetCachedTenantDB(DBURL, ctx)

	if err != nil {
		return nil, errors.New("failed to get cached tenant db to upsert teams: " + err.Error())
	}

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return nil, errors.New("failed to begin transaction to upsert teams: " + err.Error())
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	args := make([]any, 0)
	values := make([]string, 0)

	for _, t := range *teamsRecordMap {
		args = append(args, t.Name, t.ID, githubOrganizationId)
		values = append(values, "(?, ?, ?)")
	}

	query := fmt.Sprintf(
		`
		INSERT INTO github_teams 
            (name, external_id, github_organization_id) 
        VALUES 
			%s 
		ON CONFLICT 
			(external_id) 
		DO UPDATE SET 
			name = excluded.name 
		RETURNING id, name, external_id, team_id;
	`, strings.Join(values, ", "),
	)

	rows, err := tx.QueryContext(ctx, query,
		args...)

	if err != nil {
		return nil, errors.New("failed to upsert github_teams: " + err.Error())
	}

	args = make([]any, 0)
	values = make([]string, 0)

	for rows.Next() {
		var res struct {
			ID         int64
			Name       string
			ExternalID int64
			TeamID     *int64
		}

		if err := rows.Scan(&res.ID, &res.Name, &res.ExternalID, &res.TeamID); err != nil {
			return nil, errors.New("failed to scan github team upsert result: " + err.Error())
		}

		teamRecord, ok := (*teamsRecordMap)[res.Name]

		if !ok {
			return nil, errors.New("failed to get a team record from map")
		}

		teamRecord.GithubTeamID = &res.ID

		(*teamsRecordMap)[res.Name] = teamRecord

		if res.TeamID == nil {
			args = append(args, res.Name, organizationId)
			values = append(values, "(?, ?)")
		}
	}

	if len(values) > 0 {
		query = fmt.Sprintf(`
			INSERT INTO teams 
				(name, organization_id) 
			VALUES %s RETURNING id, name;`,
			strings.Join(values, ", "))

		rows, err = tx.QueryContext(ctx, query, args...)

		if err != nil {
			return nil, errors.New("failed to upsert teams: " + err.Error())
		}

		for rows.Next() {
			var res struct {
				ID   int64
				Name string
			}

			if err := rows.Scan(&res.ID, &res.Name); err != nil {
				return nil, errors.New("failed to scan team upsert result: " + err.Error())
			}

			teamRecord, ok := (*teamsRecordMap)[res.Name]

			if !ok {
				return nil, errors.New("failed to get a teamRecord from map")
			}

			teamRecord.TeamID = &res.ID
			(*teamsRecordMap)[res.Name] = teamRecord
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return teamsRecordMap, nil
}

type MemberRecord struct {
	ID    *int64
	Login *string
	Name  *string
	Email *string
	// ID of a record after insertion to github_members table
	GithubMemberId *int64
	// ID of a record after insertion to members table
	MemberID *int64
	Teams    []struct {
		Name   *string
		TeamID *int64
	}
}

type MembersRecordMap map[string]MemberRecord

func (ta *TenantActivities) UpsertGithubMembers(
	ctx context.Context,
	DBURL string,
	membersMap MembersRecordMap,
	teamsRecordMap TeamsRecordMap,
) (res *MembersRecordMap, err error) {
	db, err := ta.GetCachedTenantDB(DBURL, ctx)

	if err != nil {
		return nil, errors.New("failed to get cached tenant db to upsert teams: " + err.Error())
	}

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return nil, errors.New("failed to begin transaction to upsert teams: " + err.Error())
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	args := make([]any, 0)
	values := make([]string, 0)

	for _, m := range membersMap {
		args = append(args, m.ID, m.Login, m.Email)
		values = append(values, "(?, ?, ?)")
	}

	query := fmt.Sprintf(`
		INSERT INTO github_members
			(external_id, username, email)
		VALUES 
			%s 
		ON CONFLICT 
			(external_id) 
		DO UPDATE SET 
			username = excluded.username, 
			email = excluded.email 
		RETURNING id, member_id, username`,
		strings.Join(values, ", "))

	rows, err := tx.QueryContext(ctx, query,
		args...)

	if err != nil {
		return nil, errors.New("failed to upsert github_members: " + err.Error())
	}

	newMembersMap := MembersRecordMap{}

	args = make([]any, 0)
	values = make([]string, 0)

	for rows.Next() {
		var id, member_id *int64
		var username *string

		if err := rows.Scan(&id, &member_id, &username); err != nil {
			return nil, errors.New("failed to scan github members upsert result: " + err.Error())
		}

		memberRecord, ok := membersMap[*username]

		if !ok {
			return nil, errors.New("failed to get a team from map")
		}

		for idx, t := range memberRecord.Teams {
			team, ok := teamsRecordMap[*t.Name]
			t.TeamID = team.TeamID

			memberRecord.Teams[idx] = t

			if !ok {
				return nil, errors.New("failed to get a team from map")
			}

			args = append(args, team.GithubTeamID, id)
			values = append(values, "(?, ?)")
		}

		if member_id == nil {
			memberRecord.GithubMemberId = id
			newMembersMap[*username] = memberRecord
		}
	}

	query = fmt.Sprintf(`
		INSERT OR IGNORE INTO github_teams__github_members
			(github_team_id, github_member_id)
		VALUES %s ;`,
		strings.Join(values, ", "))

	_, err = tx.QueryContext(ctx, query,
		args...)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &newMembersMap, nil
}

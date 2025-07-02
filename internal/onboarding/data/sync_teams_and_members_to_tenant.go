package data

import (
	"context"
	"database/sql"
	"fmt"
)

func SyncTeamsAndMembersToTenant(
	ctx context.Context,
	teamWithMembers TeamWithMembers,
	dbUrl string,
	githubOrganizationId int64,
	organizationId int64,
	db *sql.DB,
) (bool, error) {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return false, err
	}

	teamsData, err := SyncTeamsToTenant(tx, ctx, teamWithMembers.Team.ID, teamWithMembers.Team.Name, organizationId, githubOrganizationId)

	if err != nil {
		return false, err
	}

	githubTeamId := teamsData.GithubTeamId
	teamId := teamsData.TeamId

	for _, member := range teamWithMembers.Members {

		_, err = tx.Exec(`
			INSERT INTO github_members
				(external_id, username, email)
			VALUES
				(?, ?, ?)
			ON CONFLICT 
				(external_id) 
			DO NOTHING;`,
			member.ID, member.Login, member.Email)

		if err != nil {
			fmt.Println("Issue while creating github member")

			_ = tx.Rollback()
			return false, err
		}

		rowRes := tx.QueryRowContext(ctx, `
			SELECT 
				id, member_id
			FROM 
				github_members 
			WHERE 
				external_id = ?;`,
			member.ID)

		var githubMemberId int64
		var memberRefId *int64

		err := rowRes.Scan(&githubMemberId, &memberRefId)

		if err != nil {
			fmt.Println("Issue while retrieving github member")
			_ = tx.Rollback()
			return false, err
		}

		_, err = tx.Exec(`
			INSERT INTO github_teams__github_members
				(github_team_id, github_member_id)
			VALUES
				(?, ?);`,
			githubTeamId, githubMemberId)

		if err != nil {
			fmt.Println("Issue creating github_teams__github_members")
			_ = tx.Rollback()
			return false, err
		}

		if memberRefId == nil {
			name := member.Name

			if name == nil {
				defaultName := "DXTA member"
				name = &defaultName
			}

			rowRes = tx.QueryRowContext(ctx, `
			INSERT INTO members
				(name, email)
			VALUES
				(?, ?)
			RETURNING id;`,
				name, member.Email)

			var memberId int64

			err = rowRes.Scan(&memberId)

			if err != nil {
				fmt.Println("Issue creating member")

				_ = tx.Rollback()
				return false, err
			}

			_, err = tx.Exec(`
				UPDATE 
					github_members 
				SET 
					member_id = ? 
				WHERE id = ?`,
				memberId, githubMemberId)

			if err != nil {
				fmt.Println("Issue while updating member_id in github member")

				_ = tx.Rollback()
				return false, err
			}

			_, err = tx.Exec(`
			INSERT INTO teams__members
				(team_id, member_id)
			VALUES
				(?, ?);`,
				teamId, memberId)

			if err != nil {
				fmt.Println("Issue creating teams__members")
				_ = tx.Rollback()
				return false, err
			}
		} else {
			_, err = tx.Exec(`
			INSERT INTO teams__members
				(team_id, member_id)
			VALUES
				(?, ?);`,
				teamId, memberRefId)

			if err != nil {
				fmt.Println("Issue creating teams__members")
				_ = tx.Rollback()
				return false, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

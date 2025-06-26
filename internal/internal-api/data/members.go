package data

import (
	"context"
	"fmt"
)

type CreateMemberResponse struct{ Id int64 }

func (d TenantDB) CreateMember(
	name string,
	email *string,
	ctx context.Context,
) (*CreateMemberResponse, error) {
	query := `
		INSERT INTO members
			(name, email)
		VALUES
			(?, ?)
		RETURNING id;`

	rows := d.DB.QueryRowContext(ctx, query, name, email)

	var newMemberId int64

	err := rows.Scan(&newMemberId)

	if err != nil {
		fmt.Printf("Could not create member with name: %s. Error: %s", name, err.Error())
		return nil, err
	}

	return &CreateMemberResponse{
		Id: newMemberId,
	}, nil

}

func (d TenantDB) AddMemberToTeam(teamId int64, memberId int64, ctx context.Context) error {
	query := `
		INSERT INTO teams__members
			(team_id, member_id)
		VALUES
			(?, ?);`

	_, err := d.DB.Exec(query, teamId, memberId)

	if err != nil {
		fmt.Printf(
			"Could not add member with id: %d to team with id %d . Error: %s",
			memberId,
			teamId,
			err.Error(),
		)
		return err
	}

	return nil
}

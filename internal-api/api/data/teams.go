package data

import (
	"context"
	"fmt"
)

type CreateTeamResponse struct{ Id int64 }

func (d TenantDB) CreateTeam(teamName string, organizationId string, ctx context.Context) (*CreateTeamResponse, error) {
	query := `
		INSERT INTO teams 
			(name, organization_id) 
		VALUES 
			(?, ?) 
		RETURNING id;`

	rows := d.DB.QueryRowContext(ctx, query, teamName, organizationId)

	var newTeamId int64

	err := rows.Scan(&newTeamId)

	if err != nil {
		fmt.Printf("Could not create team with name: %s for organization with id: %s. Error: %s", teamName, organizationId, err.Error())
		return nil, err
	}

	return &CreateTeamResponse{
		Id: newTeamId,
	}, nil
}

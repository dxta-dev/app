package data

import (
	"context"
	"fmt"
)

type CreateTeamResponse struct{ Id int64 }

func (d TenantDB) CreateTeam(teamName string, organizationId int64, ctx context.Context) (*CreateTeamResponse, error) {
	query := `
		INSERT INTO teams 
			(name, organization_id) 
		VALUES 
			(?, ?) 
		RETURNING id;`

	var newTeamId int64

	if err := d.DB.QueryRowContext(ctx, query, teamName, organizationId).Scan(&newTeamId); err != nil {
		fmt.Printf("Could not create team with name: %s for organization with id: %d. Error: %s", teamName, organizationId, err.Error())
		return nil, err
	}

	return &CreateTeamResponse{
		Id: newTeamId,
	}, nil
}

package data

import "context"

type CreateCustomTeamRes struct{ Id string }

func (d TenantDB) CreateCustomTeam(teamName string, organizationId string, ctx context.Context) (*CreateCustomTeamRes, error) {
	query := `
		INSERT INTO teams 
			(name, organization_id) 
		VALUES 
			(?, ?) 
		RETURNING id;`

	rows := d.DB.QueryRowContext(ctx, query, teamName, organizationId)

	var newTeamId string

	err := rows.Scan(&newTeamId)

	if err != nil {
		return nil, err
	}

	return &CreateCustomTeamRes{
		Id: newTeamId,
	}, nil
}

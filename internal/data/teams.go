package data

import (
	"context"
)

type Team struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (d DB) GetTeams(ctx context.Context) ([]Team, error) {
	query := "SELECT id, name FROM tenant_teams;"

	rows, err := d.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var teams = make([]Team, 0)
	for rows.Next() {
		var team Team

		if err := rows.Scan(
			&team.Id,
			&team.Name,
		); err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}

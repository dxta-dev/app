package api

import (
	"context"
	"database/sql"
)

type Team struct {
	Id   int64
	Name string
}

func GetTeams(db *sql.DB, ctx context.Context) ([]Team, error) {
	query := "SELECT id, name FROM tenant_teams;"

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var teams []Team = make([]Team, 0)
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

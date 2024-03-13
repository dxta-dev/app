package data

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"

	_ "github.com/libsql/libsql-client-go/libsql"
)

type Team struct {
	Id   int64
	Name string
}

type TeamSlice []Team

func (s *Store) GetTeams() (TeamSlice, error) {
	db, err := sql.Open("libsql", s.DbUrl)

	query := `SELECT id, name FROM tenant_teams`

	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var teams TeamSlice

	for rows.Next() {
		var team Team

		if err := rows.Scan(
			&team.Id,
			&team.Name,
		); err != nil {
			log.Fatal(err)
		}

		teams = append(teams, team)
	}

	return teams, nil
}
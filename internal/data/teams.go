package data

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"

	_ "github.com/tursodatabase/go-libsql"
)

type Team struct {
	Id   int64
	Name string
}

type TeamSlice []Team

func (s *Store) GetTeams() (TeamSlice, error) {
	db, err := sql.Open(s.DriverName, s.DbUrl)

	query := `SELECT id, name FROM tenant_teams`

	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.QueryContext(s.Context, query)

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

func (s *Store) GetTeamMembers(team *int64) ([]int64, error) {
	if team == nil {
		return []int64{}, nil
	}

	db, err := sql.Open(s.DriverName, s.DbUrl)

	query := `SELECT member FROM tenant_team_members where team = ?`

	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.QueryContext(s.Context, query, team)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var teamMembers []int64

	for rows.Next() {
		var member int64

		if err := rows.Scan(
			&member,
		); err != nil {
			log.Fatal(err)
		}

		teamMembers = append(teamMembers, member)
	}

	return teamMembers, nil
}

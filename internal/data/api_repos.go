package data

import "database/sql"

type Repo struct {
	ID           int
	Organization string
	Repository   string
	DBURL        string
}

func GetRepos(db *sql.DB) ([]Repo, error) {
	query := `SELECT id, organization, repository, db_url FROM repos`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []Repo
	for rows.Next() {
		var repo Repo
		err := rows.Scan(&repo.ID, &repo.Organization, &repo.Repository, &repo.DBURL)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}

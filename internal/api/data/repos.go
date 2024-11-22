package data

import (
	"context"
	"database/sql"
)

type Repo struct {
	Organization       string `json:"organization"`
	Repository         string `json:"repository"`
	ProjectName        string `json:"projectName"`
	ProjectDescription string `json:"projectDescription"`
}

func GetReposDbUrl(ctx context.Context, db *sql.DB, org string, repo string) (string, error) {
	query := "SELECT db_url FROM repos WHERE organization = ? AND repository = ?"

	row := db.QueryRowContext(ctx, query, org, repo)

	var dbUrl string

	err := row.Scan(&dbUrl)

	if err != nil {
		return "", err
	}

	return dbUrl, nil
}

func GetRepos(ctx context.Context, db *sql.DB) ([]Repo, error) {
	query := `SELECT organization, repository, project_name, project_description FROM repos`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []Repo
	for rows.Next() {
		var repo Repo
		err := rows.Scan(&repo.Organization, &repo.Repository, &repo.ProjectName, &repo.ProjectDescription)
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

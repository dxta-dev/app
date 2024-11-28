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
	query := `
		SELECT t.db_url
		FROM repos AS r
		JOIN tenants AS t
		ON t.id = r.tenant_id
		WHERE r.organization = ?
		AND r.repository = ?;
	`

	row := db.QueryRowContext(ctx, query, org, repo)

	var dbUrl string

	err := row.Scan(&dbUrl)

	if err != nil {
		return "", err
	}

	return dbUrl, nil
}

func GetRepos(ctx context.Context, db *sql.DB) ([]Repo, error) {
	query := `
		SELECT
			r.organization,
			r.repository,
			r.project_name,
			COALESCE(ci.description, '') AS project_description
		FROM repos AS r
		JOIN tenants AS t
		ON t.id = r.tenant_id
		LEFT JOIN company_info AS ci
		ON t.id = ci.tenant_id;
	`
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

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
	InternalTeam       int32  `json:"teamId"`
}
type TenantRepo struct {
	DbUrl        string
	Organization string
	Repository   string
}

func GetTenantRepo(ctx context.Context, db *sql.DB, org string, repo string) (TenantRepo, error) {
	query := `
		SELECT t.db_url, r.organization, r.repository
		FROM repos AS r
		JOIN tenants AS t
		ON t.id = r.tenant_id
		WHERE LOWER(r.organization) = LOWER(?)
		AND LOWER(r.repository) = LOWER(?);
	`
	row := db.QueryRowContext(ctx, query, org, repo)

	var tenantRepo TenantRepo

	err := row.Scan(&tenantRepo.DbUrl, &tenantRepo.Organization, &tenantRepo.Repository)

	if err != nil {
		return TenantRepo{}, err
	}

	return tenantRepo, nil
}

func GetRepos(ctx context.Context, db *sql.DB) ([]Repo, error) {
	query := `
		SELECT
			r.organization,
			r.repository,
			r.project_name,
			COALESCE(ci.description, '') AS project_description,
			r.internal_team
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

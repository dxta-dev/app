package onboarding

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"sync"

	internal_api_data "github.com/dxta-dev/app/internal/internal-api/data"
)

type CreateTenantConfig struct {
	TursoApiURL                 string
	TursoOrganizationSlug       string
	TursoDBGroupName            string
	TursoAuthToken              string
	TenantSeedDBURL             string
	OrganizationsTenantMapDBURL string
}

func LoadCreateTenantConfig() (*CreateTenantConfig, error) {
	var tursoAuthToken, tenantSeedDBURL, organizationsTenantMapDBURL, tursoApiUrl, tursoOrganizationSlug, tursoDBGroupName string

	if tursoAuthToken = os.Getenv("TURSO_AUTH_TOKEN"); tursoAuthToken == "" {
		return nil, errors.New("turso auth token not defined")
	}

	if tenantSeedDBURL = os.Getenv("TENANT_SEED_DB_NAME"); tenantSeedDBURL == "" {
		return nil, errors.New("seed db url not defined")
	}

	if organizationsTenantMapDBURL = os.Getenv("ORGANIZATIONS_TENANT_MAP_DB_URL"); organizationsTenantMapDBURL == "" {
		return nil, errors.New("organizations tenant map db url not defined")
	}

	if tursoApiUrl = os.Getenv("TURSO_API_URL"); tursoApiUrl == "" {
		return nil, errors.New("turso api url not defined")
	}

	if tursoOrganizationSlug = os.Getenv("TURSO_ORGANIZATION_SLUG"); tursoOrganizationSlug == "" {
		return nil, errors.New("turso organization slug not defined")
	}

	if tursoDBGroupName = os.Getenv("TURSO_DB_GROUP_NAME"); tursoDBGroupName == "" {
		return nil, errors.New("turso db group name not defined")
	}

	return &CreateTenantConfig{
		TursoApiURL:                 tursoApiUrl,
		TursoOrganizationSlug:       tursoOrganizationSlug,
		TursoDBGroupName:            tursoDBGroupName,
		TursoAuthToken:              tursoAuthToken,
		TenantSeedDBURL:             tenantSeedDBURL,
		OrganizationsTenantMapDBURL: organizationsTenantMapDBURL,
	}, nil

}

func GetCachedTenantDB(store *sync.Map, dbUrl string, ctx context.Context) (*sql.DB, error) {
	db, ok := store.Load(dbUrl)

	if !ok {
		tenantDB, err := internal_api_data.NewDB(dbUrl, ctx)

		if err != nil {
			return nil, errors.New("failed to create tenant db connection: " + err.Error())
		}

		db = tenantDB.DB
		store.Store(dbUrl, db)
	}

	return db.(*sql.DB), nil
}

func GetDB(ctx context.Context, DBURL string) (*sql.DB, error) {
	db, err := internal_api_data.NewDB(DBURL, ctx)

	if err != nil {
		return nil, errors.New("failed to create tenant db connection: " + err.Error())
	}

	return db.DB, nil
}

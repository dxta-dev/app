package onboarding

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"sync"

	"github.com/dxta-dev/app/internal/otel"
)

var tenantDBConnections = sync.Map{}

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

type DB struct {
	DB *sql.DB
}

func NewDB(DBURL string, ctx context.Context) (DB, error) {
	driverName := otel.GetDriverName()
	devToken := os.Getenv("DXTA_DEV_GROUP_TOKEN")

	if devToken == "" {
		return DB{}, errors.New("no dev group token provided")
	}

	tenantDB, err := sql.Open(
		driverName,
		DBURL+"?authToken="+devToken,
	)

	if err != nil {
		return DB{}, errors.New("failed to open db connection " + err.Error())
	}

	if err := tenantDB.PingContext(ctx); err != nil {
		return DB{}, errors.New("failed to verify db connection " + err.Error())
	}

	return DB{
		DB: tenantDB,
	}, nil
}

func GetCachedTenantDB(DBURL string, ctx context.Context) (*sql.DB, error) {
	if cachedDB, ok := tenantDBConnections.Load(DBURL); ok {
		return cachedDB.(*sql.DB), nil
	}

	db, err := NewDB(DBURL, ctx)

	if err != nil {
		return nil, errors.New("failed to create tenant db connection: " + err.Error())
	}

	tenantDBConnections.Store(DBURL, db.DB)

	return db.DB, nil
}

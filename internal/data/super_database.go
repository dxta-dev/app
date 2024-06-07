package data

import (
	"context"
	"database/sql"

	"github.com/dxta-dev/app/internal/otel"
)

type Tenant struct {
	Subdomain string
	DatabaseUrl string
}

func GetTenantBySubdomain(ctx context.Context, superDatabaseUrl string, subdomain string) (Tenant, error) {
	db, err := sql.Open(otel.GetDriverName(), superDatabaseUrl)

	if err != nil {
		return Tenant{}, err
	}

	defer db.Close()

	query := `
		SELECT
			db_url
		FROM tenants
		WHERE subdomain = ?
	`

	var tenant Tenant

	tenant.Subdomain = subdomain


	err = db.QueryRowContext(ctx, query).Scan(
		&tenant.DatabaseUrl,
	)

	if (err != nil) {
		return Tenant{}, err
	}

	return tenant, nil
}

func GetTenants(ctx context.Context, superDatabaseUrl string) ([]Tenant, error) {
	db, err := sql.Open(otel.GetDriverName(), superDatabaseUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	query := `
		SELECT
			subdomain,
			db_url
		FROM tenants
	`

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tenants []Tenant

	for rows.Next() {
		var tenant Tenant
		err := rows.Scan(&tenant.Subdomain, &tenant.DatabaseUrl)
		if err != nil {
			return nil, err
		}
		tenant.Name = tenant.Subdomain
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

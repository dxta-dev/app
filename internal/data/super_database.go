package data

import (
	"context"
	"database/sql"

	"github.com/dxta-dev/app/internal/config"
	"github.com/dxta-dev/app/internal/otel"
)

func getTenants(ctx context.Context, superDatabaseUrl string) ([]config.Tenant, error) {
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

	var tenants []config.Tenant

	for rows.Next() {
		var tenant config.Tenant
		err := rows.Scan(&tenant.Subdomain, &tenant.DatabaseUrl)
		if err != nil {
			return nil, err
		}
		tenant.DatabaseType = config.LibSQL
		tenant.Name = tenant.Subdomain
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

package data

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/dxta-dev/app/internal/otel"
)

type TenantDB struct {
	DB *sql.DB
}

func NewTenantDB(dbUrl string, ctx context.Context) (TenantDB, error) {
	driverName := otel.GetDriverName()
	devToken := os.Getenv("DXTA_DEV_GROUP_TOKEN")

	tenantDB, err := sql.Open(
		driverName,
		dbUrl+"?authToken="+devToken,
	)

	if err != nil {
		fmt.Printf(
			"Issue while opening tenant database connection. DBUrl: %s Error: %s",
			dbUrl,
			err.Error(),
		)
		return TenantDB{}, err
	}

	if err := tenantDB.PingContext(ctx); err != nil {
		return TenantDB{}, err
	}

	return TenantDB{
		DB: tenantDB,
	}, nil
}

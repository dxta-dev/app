package data

import (
	"context"
	"database/sql"
	"errors"
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
		return TenantDB{}, errors.New("failed to open tenant db connection " + err.Error())
	}

	if err := tenantDB.PingContext(ctx); err != nil {
		return TenantDB{}, errors.New("failed to verify tenant db connection " + err.Error())
	}

	return TenantDB{
		DB: tenantDB,
	}, nil
}

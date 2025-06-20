package data

import (
	"database/sql"
	"os"

	"github.com/dxta-dev/app/internal/otel"
)

type TenantDB struct {
	DB *sql.DB
}

func NewTenantDB(dbUrl string) (TenantDB, error) {
	driverName := otel.GetDriverName()
	devToken := os.Getenv("DXTA_DEV_GROUP_TOKEN")

	tenantDB, err := sql.Open(
		driverName,
		dbUrl+"?authToken="+devToken,
	)

	if err != nil {
		return TenantDB{}, err
	}

	return TenantDB{
		DB: tenantDB,
	}, nil
}

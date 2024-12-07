package data

import (
	"context"
	"database/sql"
	"github.com/dxta-dev/app/internal/otel"
	"os"
	"strings"
	"sync"
)

type DB struct {
	db *sql.DB
}

var dbPool sync.Map

func (database DB) get(ctx context.Context, tenantRepo TenantRepo) (*sql.DB, error) {
	cacheKey := strings.ToLower(tenantRepo.Organization + "/" + tenantRepo.Repository)

	if dbInterface, ok := dbPool.Load(cacheKey); ok {
		return dbInterface.(*sql.DB), nil
	}

	driverName := otel.GetDriverName()

	fullDbUrl := tenantRepo.DbUrl + "?authToken=" + os.Getenv("DXTA_DEV_GROUP_TOKEN")

	db, err := sql.Open(driverName, fullDbUrl)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	dbPool.Store(cacheKey, db)

	return db, nil
}

package data

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"sync"

	"github.com/dxta-dev/app/internal/otel"
)

type DB struct {
	db *sql.DB
}

var dbPool sync.Map

func NewDB(ctx context.Context, tenantRepo TenantRepo) (DB, error) {
	cacheKey := strings.ToLower(tenantRepo.Organization + "/" + tenantRepo.Repository)

	if dbInterface, ok := dbPool.Load(cacheKey); ok {
		return DB{
			db: dbInterface.(*sql.DB),
		}, nil
	}

	driverName := otel.GetDriverName()

	fullDbUrl := tenantRepo.DbUrl + "?authToken=" + os.Getenv("DXTA_DEV_GROUP_TOKEN")

	db, err := sql.Open(driverName, fullDbUrl)
	if err != nil {
		return DB{}, err
	}

	if err := db.PingContext(ctx); err != nil {
		return DB{}, err
	}

	dbPool.Store(cacheKey, db)

	return DB{
		db: db,
	}, nil
}

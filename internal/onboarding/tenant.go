package onboarding

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	internal_api_data "github.com/dxta-dev/app/internal/internal-api/data"
)

var tenantDBConnections = sync.Map{}

func GetCachedTenantDB(DBURL string, ctx context.Context) (*sql.DB, error) {
	if cachedDB, ok := tenantDBConnections.Load(DBURL); ok {
		return cachedDB.(*sql.DB), nil
	}

	db, err := internal_api_data.NewDB(DBURL, ctx)

	if err != nil {
		return nil, errors.New("failed to create tenant db connection: " + err.Error())
	}

	tenantDBConnections.Store(DBURL, db.DB)

	return db.DB, nil
}

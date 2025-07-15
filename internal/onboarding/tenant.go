package onboarding

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	internal_api_data "github.com/dxta-dev/app/internal/internal-api/data"
)

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

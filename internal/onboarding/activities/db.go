package activities

import (
	"context"
	"database/sql"
	"sync"

	internal_api_data "github.com/dxta-dev/app/internal/internal-api/data"
)

func GetCachedTenantDB(store *sync.Map, dbUrl string, ctx context.Context) (*sql.DB, error) {
	db, ok := store.Load(dbUrl)

	if !ok {
		tenantDB, err := internal_api_data.NewTenantDB(dbUrl, ctx)
		db = tenantDB.DB

		if err != nil {
			return nil, err
		}

		store.Store(dbUrl, db)
	}

	return db.(*sql.DB), nil
}

type DBActivities struct {
	Connections       sync.Map
	GetCachedTenantDB func(store *sync.Map, dbUrl string, ctx context.Context) (*sql.DB, error)
}

func InitDBActivities() *DBActivities {
	return &DBActivities{
		Connections:       sync.Map{},
		GetCachedTenantDB: GetCachedTenantDB,
	}
}

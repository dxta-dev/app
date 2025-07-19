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

type QueryKeys interface {
	AggregatedStatisticsKey | AggregatedValuesKey | CycleTimeStatisticsKey
}

type Query[T QueryKeys] struct {
	value string
	_     T
}

func NewQuery[T QueryKeys](query string) Query[T] {
	return Query[T]{
		value: query,
	}
}

func (q *Query[T]) Get() string {
	return q.value
}

type Executable[T QueryKeys] interface {
	QueryKeys
	Execute(ctx context.Context,
		db *DB,
		q Query[T],
		org string,
		repo string,
		weeks []string,
		team *int64,
	) (any, error)
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

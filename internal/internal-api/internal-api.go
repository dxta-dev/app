package api

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/dxta-dev/app/internal/internal-api/data"
	"github.com/dxta-dev/app/internal/otel"
	_ "github.com/libsql/libsql-client-go/libsql"
)

type State struct {
	DB data.DB
}

type TenantDBData struct {
	DBUrl string
}

var tenantDBURLcache sync.Map

func GetTenantDBUrlByAuthId(ctx context.Context, authID string) (TenantDBData, error) {
	if cached, ok := tenantDBURLcache.Load(authID); ok {
		return TenantDBData{DBUrl: cached.(string)}, nil
	}

	driverName := otel.GetDriverName()
	tenantOrganizationMapDBUrl := os.Getenv("TENANT_ORG_MAPPING_URL")
	devToken := os.Getenv("DXTA_DEV_GROUP_TOKEN")

	tenantOrganizationMapDB, err := sql.Open(
		driverName,
		tenantOrganizationMapDBUrl+"?authToken="+devToken,
	)

	if err != nil {
		fmt.Printf(
			"Issue while opening organizations-tenant-map database connection. Error: %s",
			err.Error(),
		)
		return TenantDBData{}, err
	}

	defer tenantOrganizationMapDB.Close()

	query := `
		SELECT db_url
		FROM tenants
		WHERE organization_id = ?;`

	var tenantData TenantDBData

	if err = tenantOrganizationMapDB.QueryRowContext(ctx, query, authID).Scan(&tenantData.DBUrl); err != nil {
		fmt.Printf(
			"Could not retrieve tenant db url for organization with id: %s. Error: %s",
			authID,
			err.Error(),
		)
		return TenantDBData{}, err
	}

	tenantDBURLcache.Store(authID, tenantData.DBUrl)

	return tenantData, nil
}

func InternalApiState(authId string, ctx context.Context) (State, error) {
	tenantData, err := GetTenantDBUrlByAuthId(ctx, authId)

	if err != nil {
		return State{}, err
	}

	tenantDB, err := data.NewDB(tenantData.DBUrl, ctx)

	if err != nil {
		return State{}, err
	}

	return State{
		DB: tenantDB,
	}, nil
}

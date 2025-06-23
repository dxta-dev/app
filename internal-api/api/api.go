package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/dxta-dev/app/internal-api/api/data"
	"github.com/dxta-dev/app/internal/otel"
	_ "github.com/libsql/libsql-client-go/libsql"
)

type State struct {
	DB data.TenantDB
}

type TenantDBData struct {
	DBUrl string
}

func GetTenantDBUrlByAuthId(ctx context.Context, authId string) (TenantDBData, error) {
	driverName := otel.GetDriverName()
	tenantOrganizationMapDBUrl := os.Getenv("TENANT_ORG_MAPPING_URL")
	devToken := os.Getenv("DXTA_DEV_GROUP_TOKEN")

	tenantOrganizationMapDB, err := sql.Open(
		driverName,
		tenantOrganizationMapDBUrl+"?authToken="+devToken,
	)

	if err != nil {
		fmt.Printf("Issue while opening organizations-tenant-map database connection. Error: %s", err.Error())
		return TenantDBData{}, err
	}

	defer tenantOrganizationMapDB.Close()

	query := `
		SELECT db_url 
		FROM tenants 
		WHERE organization_id = ?;`

	var tenantData TenantDBData

	if err = tenantOrganizationMapDB.QueryRowContext(ctx, query, authId).Scan(&tenantData.DBUrl); err != nil {
		fmt.Printf("Could not retrieve tenant db url for organization with id: %s. Error: %s", authId, err.Error())
		return TenantDBData{}, err
	}

	return tenantData, nil
}

func InternalApiState(dbUrl string, r *http.Request) (State, error) {
	tenantDB, err := data.NewTenantDB(dbUrl)

	if err != nil {
		return State{}, err
	}

	return State{
		DB: tenantDB,
	}, nil
}

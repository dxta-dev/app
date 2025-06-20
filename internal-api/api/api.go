package api

import (
	"context"
	"database/sql"
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
	dbUrl string
}

func getTenantDBUrl(ctx context.Context, organizationId string) (TenantDBData, error) {
	driverName := otel.GetDriverName()
	tenantOrganizationMapDBUrl := os.Getenv("TENANT_ORG_MAPPING_URL")
	devToken := os.Getenv("DXTA_DEV_GROUP_TOKEN")

	tenantOrganizationMapDB, err := sql.Open(
		driverName,
		tenantOrganizationMapDBUrl+"?authToken="+devToken,
	)

	if err != nil {
		return TenantDBData{}, err
	}

	defer tenantOrganizationMapDB.Close()

	query := `
		SELECT db_url 
		FROM tenants 
		WHERE organization_id = ?`
	row := tenantOrganizationMapDB.QueryRowContext(ctx, query, organizationId)

	var tenantData TenantDBData

	err = row.Scan(&tenantData.dbUrl)

	if err != nil {
		return TenantDBData{}, err
	}

	return tenantData, nil
}

func PlatformApiState(r *http.Request, organizationId string) (State, error) {
	ctx := r.Context()
	tenantData, err := getTenantDBUrl(ctx, organizationId)

	if err != nil {
		return State{}, err
	}

	tenantDB, err := data.NewTenantDB(tenantData.dbUrl)

	if err != nil {
		return State{}, err
	}

	return State{
		DB: tenantDB,
	}, nil
}

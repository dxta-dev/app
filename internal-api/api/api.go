package api

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/dxta-dev/app/internal-api/api/data"
	"github.com/dxta-dev/app/internal/otel"
)

type State struct {
	DB data.TenantDB
}

func PlatformApiState(r *http.Request, organizationId string) (State, error) {
	driverName := otel.GetDriverName()
	tenantOrganizationMapDBUrl := os.Getenv("TENANT_ORG_MAPPING_URL")
	devToken := os.Getenv("DXTA_DEV_GROUP_TOKEN")

	tenantOrganizationMapDB, err := sql.Open(
		driverName,
		tenantOrganizationMapDBUrl+"?authToken="+devToken,
	)

	if err != nil {
		return State{}, err
	}

	defer tenantOrganizationMapDB.Close()

	ctx := r.Context()

	query := `
		SELECT db_url 
		FROM tenants 
		WHERE organization_id = ?`
	row := tenantOrganizationMapDB.QueryRowContext(ctx, query, organizationId)

	var tenantData struct{ DbUrl string }

	err = row.Scan(&tenantData.DbUrl)

	if err != nil {
		return State{}, err
	}

	tenantDB, err := data.NewTenantDB(tenantData.DbUrl)

	if err != nil {
		return State{}, err
	}

	return State{
		DB: tenantDB,
	}, nil
}

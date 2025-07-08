package activity

import (
	"context"
	"errors"
	"fmt"

	"github.com/dxta-dev/app/internal/onboarding"
)

func (ta *TenantActivities) UpsertTenantDBInfo(
	ctx context.Context,
	DBName string,
	DBHostName string,
	DBDomainName string,
) (bool, error) {
	dbUrl := fmt.Sprintf("libsql://%s", DBHostName)
	db, err := onboarding.GetCachedTenantDB(ta.DBConnections, dbUrl, ctx)

	if err != nil {
		return false, err
	}

	_, err = db.QueryContext(ctx, `
		INSERT INTO settings 
			(tenant_name, tenant_domain) 
		VALUES 
			(?, ?);`,
		DBName, DBDomainName,
	)

	if err != nil {
		return false, errors.New("Failed to upsert tenant db info: " + err.Error())
	}

	return true, nil
}

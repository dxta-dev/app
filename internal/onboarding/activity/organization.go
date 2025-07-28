package activity

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dxta-dev/app/internal/onboarding"
)

type TenantActivities struct {
	GetCachedTenantDB func(dbUrl string, ctx context.Context) (*sql.DB, error)
}

func NewTenantActivities() *TenantActivities {
	return &TenantActivities{GetCachedTenantDB: onboarding.GetCachedTenantDB}
}

func (ta *TenantActivities) GetOrganizationIDByAuthID(ctx context.Context, authID string, DBURL string) (int64, error) {
	db, err := ta.GetCachedTenantDB(DBURL, ctx)

	if err != nil {
		return 0, err
	}

	var organizationId int64

	if err = db.QueryRowContext(ctx, `
		SELECT 
			id 
		FROM 
			organizations 
		WHERE 
			auth_id = ?;`,
		authID).Scan(&organizationId); err != nil {
		return 0, errors.New("failed to retrieve organization: " + err.Error())
	}

	return organizationId, nil
}

package activity

import (
	"context"
	"errors"
	"sync"

	"github.com/dxta-dev/app/internal/onboarding"
)

type TenantActivities struct {
	DBConnections *sync.Map
}

func NewTenantActivities(DBConnections *sync.Map) *TenantActivities {
	return &TenantActivities{DBConnections}
}

func (ta *TenantActivities) GetOrganizationIDByAuthID(ctx context.Context, authID string, DBURL string) (int64, error) {
	db, err := onboarding.GetCachedTenantDB(ta.DBConnections, DBURL, ctx)

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

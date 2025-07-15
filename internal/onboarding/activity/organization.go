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
		return 0, errors.New("failed to get cached tenant DB: " + err.Error())
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

func (ta *TenantActivities) CreateOrganization(
	ctx context.Context,
	organizationName string,
	authID string,
	DBURL string,
) (bool, error) {

	db, err := onboarding.GetCachedTenantDB(ta.DBConnections, DBURL, ctx)

	if err != nil {
		return false, errors.New("failed to get cached tenant DB: " + err.Error())
	}

	_, err = db.QueryContext(ctx, `
		INSERT INTO organizations 
			(name, auth_id) 
		VALUES 
			(?, ?);`,
		organizationName, authID)

	if err != nil {
		return false, errors.New("failed to create organization: " + err.Error())
	}

	return true, nil
}

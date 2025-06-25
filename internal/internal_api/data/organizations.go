package data

import (
	"context"
	"fmt"
)

func (d TenantDB) GetOrganizationIdByAuthId(authId string, ctx context.Context) (int64, error) {
	query := `
	SELECT id
	FROM organizations
	WHERE auth_id = ?;`

	var organizationId int64

	if err := d.DB.QueryRowContext(ctx, query, authId).Scan(&organizationId); err != nil {
		fmt.Printf("Could not retrieve organization with auth id: %s", authId)
		return 0, err
	}

	return organizationId, nil
}

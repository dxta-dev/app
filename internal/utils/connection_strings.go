package utils

import (
	"fmt"
	"os"
)

func GetTenantDatabaseUrl(tenant string) string {
	return fmt.Sprintf(os.Getenv("TENANT_DATABASE_URL_TEMPLATE"), tenant, os.Getenv("TENANT_DATABASE_AUTH_TOKEN"))
}

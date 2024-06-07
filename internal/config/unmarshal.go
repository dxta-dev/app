package config

import (
	"sort"

	"github.com/knadh/koanf/v2"
)

type tenant struct {
	subdomain        string
	databaseFilePath string
	databaseUrl      string
}

type config struct {
	port             int
	tenants          map[string]tenant
	superDatabaseUrl string
}

func unmarshal(k *koanf.Koanf) (config, error) {
	out := config{}

	out.port = k.Int("port")
	out.superDatabaseUrl = k.String("super_database_url")

	out.tenants = make(map[string]tenant)

	tenantKeys := k.MapKeys("tenants")

	sort.Strings(tenantKeys)

	for _, key := range tenantKeys {
		tenant := tenant{}

		tenant.subdomain = k.String("tenants." + key + ".subdomain")
		tenant.databaseFilePath = k.String("tenants." + key + ".database_file_path")
		tenant.databaseUrl = k.String("tenants." + key + ".database_url")

		out.tenants[key] = tenant
	}

	return out, nil
}

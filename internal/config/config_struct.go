package config

import "strings"

type tenant struct {
	subdomain        *string `koanf:"subdomain"`
	databaseFilePath *string `koanf:"database_file_path"`
	databaseUrl      *string `koanf:"database_url"`
}

type config struct {
	tenants          map[string]tenant `koanf:"tenants"`
	superDatabaseUrl *string           `koanf:"super_database_url"`
}

type DatabaseType int32

const (
	SQLite DatabaseType = iota
	LibSQL
)

type Tenant struct {
	Name              string
	Subdomain         string
	DatabaseType      DatabaseType
	LocalDatabasePath string
	DatabaseUrl       string
}

type Config struct {
	IsMultiTenant          bool
	IsSuperDatabaseEnabled bool
	SuperDatabaseUrl       string
	tenants                []Tenant
}

func getConfig(config *config) Config {

	out := Config{}

	if config.superDatabaseUrl != nil && *config.superDatabaseUrl != "" {
		out.SuperDatabaseUrl = *config.superDatabaseUrl
		out.IsSuperDatabaseEnabled = true
		out.IsMultiTenant = false
		return out
	}

	if config.tenants != nil && len(config.tenants) > 1 {
		out.IsMultiTenant = true
	}

	if config.tenants != nil {
		for key, t := range config.tenants {
			tenant := Tenant{
				Name: key,
			}

			if t.subdomain != nil && *t.subdomain != "" {
				tenant.Subdomain = *t.subdomain
			} else {
				tenant.Subdomain = strings.ToLower(key)
			}

			if t.databaseUrl != nil && *t.databaseUrl != "" {
				tenant.DatabaseUrl = *t.databaseUrl
				tenant.DatabaseType = LibSQL
			}

			if t.databaseFilePath != nil && *t.databaseFilePath != "" {
				tenant.LocalDatabasePath = *t.databaseFilePath
				tenant.DatabaseType = SQLite
			}

			out.tenants = append(out.tenants, tenant)
		}
	}

	return out
}

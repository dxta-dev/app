package config

import (
	"errors"
	"strings"
)

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
	Tenants                []Tenant
}

func getConfig(config *config) (Config, error) {

	out := Config{}

	if config.superDatabaseUrl != nil && *config.superDatabaseUrl != "" {
		out.SuperDatabaseUrl = *config.superDatabaseUrl
		out.IsSuperDatabaseEnabled = true
		out.IsMultiTenant = true
		return out, nil
	}

	if config.tenants != nil && len(config.tenants) > 1 {
		out.IsMultiTenant = true
	}

	if (config.superDatabaseUrl == nil || *config.superDatabaseUrl != "") && (config.tenants == nil || len(config.tenants) == 0) {
		return Config{}, errors.New("both super database url and tenants cannot be empty")
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

			if (t.databaseUrl == nil || *t.databaseUrl == "") && (t.databaseFilePath == nil || *t.databaseFilePath == "")  {
				return Config{}, errors.New("both database url and file path cannot be empty")
			}

			out.Tenants = append(out.Tenants, tenant)
		}
	}

	return out, nil
}

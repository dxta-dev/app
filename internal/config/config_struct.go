package config

import (
	"context"
	"errors"
	"strings"

	"github.com/dxta-dev/app/internal/data"
)


type Tenant struct {
	Name              string
	Subdomain         string
	DatabaseUrl       string
}

type Config struct {
	IsMultiTenant          bool
	IsSuperDatabaseEnabled bool
	SuperDatabaseUrl       string
	Tenants                []Tenant
}

func getConfig(ctx context.Context, config *config) (Config, error) {

	out := Config{}

	if config.superDatabaseUrl != "" {
		out.SuperDatabaseUrl = config.superDatabaseUrl
		out.IsSuperDatabaseEnabled = true
		out.IsMultiTenant = true
		ts, err := data.GetTenants(ctx, config.superDatabaseUrl)

		if err != nil {
			return Config{}, err
		}

		for _, t := range ts {
			tenant := Tenant {
				Name: t.Name,
				Subdomain: t.Subdomain,
				DatabaseUrl: t.DatabaseUrl,
			}

			out.Tenants = append(out.Tenants, tenant)
		}


		return out, nil
	}

	if config.tenants != nil && len(config.tenants) > 1 {
		out.IsMultiTenant = true
	}

	if (config.superDatabaseUrl != "") && (config.tenants == nil || len(config.tenants) == 0) {
		return Config{}, errors.New("both super database url and tenants cannot be empty")
	}

	if config.tenants != nil {
		for key, t := range config.tenants {
			tenant := Tenant{
				Name: key,
			}

			if t.subdomain != "" {
				tenant.Subdomain = t.subdomain
			} else {
				tenant.Subdomain = strings.ToLower(key)
			}

			if t.databaseUrl != "" {
				tenant.DatabaseUrl = t.databaseUrl
			}

			if t.databaseFilePath != "" {
				tenant.DatabaseUrl = t.databaseFilePath
			}

			if t.databaseUrl == "" && t.databaseFilePath == "" {
				return Config{}, errors.New("both database url and file path cannot be empty")
			}

			out.Tenants = append(out.Tenants, tenant)
		}
	}

	return out, nil
}

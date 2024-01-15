package utils

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)


type Tenant struct {
	Name string `toml:"name"`
	SubdomainName string `toml:"subdomain"`
	DatabaseName string `toml:"database_name"`
	DatabaseUrl *string `toml:"database_url"`
}

type Config struct {
	IsMultiTenant bool
	ShouldUseSuperDatabase bool
	SuperDatabaseUrl *string
	TenantDatabaseUrlTemplate *string
	Tenants map[string]Tenant
}

type TomlConfig struct {
	TenantDatabaseUrlTemplate *string `toml:"tenant_database_url_template"`
	TenantDatabaseGroupAuth *string `toml:"tenant_database_group_auth"`
	SuperDatabaseUrl *string `toml:"super_database_url"`
	Tenants map[string]Tenant `toml:"tenants"`
}

func ValidateConfig(config *TomlConfig) (*Config, error) {
	var superDatabaseUrl *string

	if config.SuperDatabaseUrl != nil {
		superDatabaseUrl = config.SuperDatabaseUrl
	}

	for key, tenant := range config.Tenants {
		if tenant.Name == "" {
			tenant.Name = key
		}
		if tenant.SubdomainName == "" {
			tenant.SubdomainName = key
		}
		if tenant.DatabaseName == "" {
			tenant.DatabaseName = key
		}

		if tenant.DatabaseUrl == nil && config.TenantDatabaseUrlTemplate != nil {
			databaseUrl := fmt.Sprintf(*config.TenantDatabaseUrlTemplate, tenant.DatabaseName)
			tenant.DatabaseUrl = &databaseUrl
		}

		if tenant.DatabaseUrl != nil && config.TenantDatabaseGroupAuth != nil {
			databaseUrl := *tenant.DatabaseUrl + "?auth_token=" + *config.TenantDatabaseGroupAuth
			tenant.DatabaseUrl = &databaseUrl
		}

		config.Tenants[key] = tenant
	}

	shouldUseSuperDatabase := false
	isMultiTenant := false


	var tenantDatabaseUrlTemplate *string
	tenantDatabaseUrlTemplate = nil

	if superDatabaseUrl != nil && config.Tenants == nil {
		shouldUseSuperDatabase = true
		isMultiTenant = true
		if (config.TenantDatabaseUrlTemplate != nil) {
			tenantDatabaseUrlTemplate = config.TenantDatabaseUrlTemplate
		}
		if (config.TenantDatabaseUrlTemplate != nil && config.TenantDatabaseGroupAuth != nil) {
			newTenantDatabaseUrlTemplate := *tenantDatabaseUrlTemplate + "?auth_token=" + *config.TenantDatabaseGroupAuth
			tenantDatabaseUrlTemplate = &newTenantDatabaseUrlTemplate
		}
	} else {
		superDatabaseUrl = nil
	}



	if len(config.Tenants) > 1 {
		isMultiTenant = true
	}

	return &Config{
		IsMultiTenant: isMultiTenant,
		ShouldUseSuperDatabase: shouldUseSuperDatabase,
		SuperDatabaseUrl: superDatabaseUrl,
		TenantDatabaseUrlTemplate: tenantDatabaseUrlTemplate,
		Tenants: config.Tenants,
	}, nil
}

func GetConfig() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if(path == "") {
		path = "config.toml"
	}
	conf, _ := LoadConfigToml(path)

	config, err := ValidateConfig(conf)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func LoadConfigToml(path string) (*TomlConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config TomlConfig
	if _, err := toml.Decode(string(data), &config); err != nil {
		return nil, err
	}

	return &config, nil
}

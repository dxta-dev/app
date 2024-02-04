package utils

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Tenant struct {
	Name          string  `toml:"name"`
	SubdomainName string  `toml:"subdomain"`
	DatabaseName  string  `toml:"database_name"`
	DatabaseUrl   *string `toml:"database_url"`
}

type Config struct {
	IsMultiTenant             bool
	ShouldUseSuperDatabase    bool
	SuperDatabaseUrl          *string
	TenantDatabaseUrlTemplate *string
	Tenants                   map[string]Tenant
}

type TomlConfig struct {
	TenantDatabaseUrlTemplate *string           `toml:"tenant_database_url_template"`
	SuperDatabaseUrl          *string           `toml:"super_database_url"`
	Tenants                   map[string]Tenant `toml:"tenants"`
}

func ValidateConfig(config *TomlConfig) (*Config, error) {
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

		if tenant.DatabaseUrl == nil {
			return nil, fmt.Errorf("utils: config needs to define either \"[tenants.%s].database_url\" or \"tenant_database_url_template\"", key)
		}

		config.Tenants[key] = tenant
	}

	configTenantsSize := len(config.Tenants)

	if configTenantsSize == 0 && config.SuperDatabaseUrl == nil {
		return nil, fmt.Errorf("utils: config needs to define either \"[tenants.*]\" or \"super_database_url\"")
	}

	if configTenantsSize > 0 {
		if config.SuperDatabaseUrl != nil {
			// TODO: logger ? after echo starts ?
			fmt.Printf("utils: using %d config \"[tenants.*]\", ignoring config \"super_database_url\"\n", configTenantsSize)
		}

		return &Config{
			IsMultiTenant:             configTenantsSize > 1,
			ShouldUseSuperDatabase:    false,
			SuperDatabaseUrl:          nil,
			TenantDatabaseUrlTemplate: config.TenantDatabaseUrlTemplate,
			Tenants:                   config.Tenants,
		}, nil
	}

	var tenantDatabaseUrlTemplate = config.TenantDatabaseUrlTemplate
	var stringIdentityTemplate = "%s"

	if tenantDatabaseUrlTemplate == nil {
		tenantDatabaseUrlTemplate = &stringIdentityTemplate
	}

	return &Config{
		IsMultiTenant:             true,
		ShouldUseSuperDatabase:    true,
		SuperDatabaseUrl:          config.SuperDatabaseUrl,
		TenantDatabaseUrlTemplate: tenantDatabaseUrlTemplate,
		Tenants:                   config.Tenants,
	}, nil
}

func GetConfigFromEnv() (*Config, error) {
	superDatabaseUrl := os.Getenv("SUPER_DATABASE_URL")
	if superDatabaseUrl == "" {
		return nil, fmt.Errorf("missing environment variable \"SUPER_DATABASE_URL\"")
	}

	groupAuthToken := os.Getenv("GROUP_AUTH_TOKEN")
	if groupAuthToken == "" {
		return nil, fmt.Errorf("missing environment variable \"GROUP_AUTH_TOKEN\"")
	}

	tenantDatabaseUrlTemplate := fmt.Sprintf("%%s?authToken=%s", groupAuthToken)

	return &Config{
		IsMultiTenant:             true,
		ShouldUseSuperDatabase:    true,
		SuperDatabaseUrl:          &superDatabaseUrl,
		TenantDatabaseUrlTemplate: &tenantDatabaseUrlTemplate,
		Tenants:                   nil,
	}, nil

}

func GetConfig() (*Config, error) {
	useEnv := os.Getenv("USE_SUPER_ENV")
	if useEnv == "true" {
		return GetConfigFromEnv()
	}
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config.toml"
	}
	conf, err := LoadConfigToml(path)

	if err != nil {
		return nil, err
	}

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

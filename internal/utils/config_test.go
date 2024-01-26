package utils

import (
	"os"
	"reflect"
	"testing"
)

func createTempFile(content string) (*os.File, error) {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		return nil, err
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		return nil, err
	}

	if err := tmpfile.Close(); err != nil {
		return nil, err
	}

	return tmpfile, nil
}

func TestValidateConfig(t *testing.T) {
	oneDatabaseUrl := "sqlite://one.db"
	superDatabaseUrl := "sqlite://super.db"

	dbTemplate := "sqlite://%s.db"

	templateDatabaseUrl := "sqlite://%s.db"

	tests := []struct {
		name      string
		config    *TomlConfig
		expected  *Config
		expectErr bool
	}{
		{
			name: "simple_config",
			config: &TomlConfig{
				SuperDatabaseUrl:          nil,
				TenantDatabaseUrlTemplate: nil,
				Tenants: map[string]Tenant{
					"tenant1": {Name: "Tenant One", SubdomainName: "one", DatabaseName: "one", DatabaseUrl: &oneDatabaseUrl},
				},
			},
			expected: &Config{
				IsMultiTenant:             false,
				ShouldUseSuperDatabase:    false,
				SuperDatabaseUrl:          nil,
				TenantDatabaseUrlTemplate: nil,
				Tenants: map[string]Tenant{
					"tenant1": {Name: "Tenant One", SubdomainName: "one", DatabaseName: "one", DatabaseUrl: &oneDatabaseUrl},
				},
			},
			expectErr: false,
		},
		{
			name: "should_set_super_database_url_to_nil",
			config: &TomlConfig{
				SuperDatabaseUrl:          &superDatabaseUrl,
				TenantDatabaseUrlTemplate: nil,
				Tenants: map[string]Tenant{
					"tenant1": {Name: "Tenant One", SubdomainName: "one", DatabaseName: "one", DatabaseUrl: &oneDatabaseUrl},
				},
			},
			expected: &Config{
				IsMultiTenant:             false,
				ShouldUseSuperDatabase:    false,
				SuperDatabaseUrl:          nil,
				TenantDatabaseUrlTemplate: nil,
				Tenants: map[string]Tenant{
					"tenant1": {Name: "Tenant One", SubdomainName: "one", DatabaseName: "one", DatabaseUrl: &oneDatabaseUrl},
				},
			},
			expectErr: false,
		},
		{
			name: "should_set_database_url_to_tenants",
			config: &TomlConfig{
				SuperDatabaseUrl:          nil,
				TenantDatabaseUrlTemplate: &dbTemplate,
				Tenants: map[string]Tenant{
					"tenant1": {Name: "Tenant One", SubdomainName: "one", DatabaseName: "one", DatabaseUrl: nil},
				},
			},
			expected: &Config{
				IsMultiTenant:             false,
				ShouldUseSuperDatabase:    false,
				SuperDatabaseUrl:          nil,
				TenantDatabaseUrlTemplate: &dbTemplate,
				Tenants: map[string]Tenant{
					"tenant1": {Name: "Tenant One", SubdomainName: "one", DatabaseName: "one", DatabaseUrl: &oneDatabaseUrl},
				},
			},
			expectErr: false,
		},
		{
			name: "should_set_default_values",
			config: &TomlConfig{
				SuperDatabaseUrl:          nil,
				TenantDatabaseUrlTemplate: nil,
				Tenants: map[string]Tenant{
					"tenant1": {Name: "", SubdomainName: "", DatabaseName: "", DatabaseUrl: &oneDatabaseUrl},
				},
			},
			expected: &Config{
				IsMultiTenant:             false,
				ShouldUseSuperDatabase:    false,
				SuperDatabaseUrl:          nil,
				TenantDatabaseUrlTemplate: nil,
				Tenants: map[string]Tenant{
					"tenant1": {Name: "tenant1", SubdomainName: "tenant1", DatabaseName: "tenant1", DatabaseUrl: &oneDatabaseUrl},
				},
			},
			expectErr: false,
		},
		{
			name: "should_use_super_database",
			config: &TomlConfig{
				SuperDatabaseUrl:          &superDatabaseUrl,
				TenantDatabaseUrlTemplate: &dbTemplate,
				Tenants:                   nil,
			},
			expected: &Config{
				IsMultiTenant:             true,
				ShouldUseSuperDatabase:    true,
				SuperDatabaseUrl:          &superDatabaseUrl,
				TenantDatabaseUrlTemplate: &templateDatabaseUrl,
				Tenants:                   nil,
			},
			expectErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := ValidateConfig(test.config)
			if err != nil && !test.expectErr {
				t.Fatalf("Expected no error, got %v", err)
			}

			if !reflect.DeepEqual(cfg, test.expected) {
				t.Errorf("\n\tExpected %+v\n\tReceived %+v", test.expected, cfg)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {

	t.Run("valid_path_and_format", func(t *testing.T) {
		content := `
super_database_url = "sqlite://super.db"
  [tenants.tenant1]
  name = "Tenant One"
  subdomain = "one"
  database_url = "sqlite://one.db"
`
		tmpfile, err := createTempFile(content)
		if err != nil {
			t.Fatal(err)
		}

		defer os.Remove(tmpfile.Name())

		superDatabaseUrl := "sqlite://super.db"
		tenant1DatabaseUrl := "sqlite://one.db"
		expected := &TomlConfig{
			SuperDatabaseUrl: &superDatabaseUrl,
			Tenants: map[string]Tenant{
				"tenant1": {Name: "Tenant One", SubdomainName: "one", DatabaseUrl: &tenant1DatabaseUrl},
			},
		}

		cfg, err := LoadConfigToml(tmpfile.Name())
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if !reflect.DeepEqual(cfg, expected) {
			t.Errorf("Expected %+v, got %+v", expected, cfg)
		}

	})

	t.Run("file_not_found", func(t *testing.T) {
		_, err := LoadConfigToml("nonexistent.toml")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("bad_format", func(t *testing.T) {
		content := `
super_database_url = "sqlite://super.db"
  [tenants.tenant1]
  name = "Tenant One"
  subdomain "one"
  database_url = sqlite://one.db
`
		tmpfile, err := createTempFile(content)
		if err != nil {
			t.Fatal(err)
		}

		defer os.Remove(tmpfile.Name())

		_, err = LoadConfigToml(tmpfile.Name())
		if err == nil {
			t.Fatalf("Somehow the config is loaded successfully")
		}
	})

	t.Run("real_life_example", func(t *testing.T) {
		content := `
super_database_url = "sqlite://super.db"

[tenants.crocoder]
  subdomain = "crocoder"
  name = "crocoder"
  database_url = "sqlite://crocoder.db"

[tenants.acme]
  database_url = "sqlite://acme.db"
`
		tmpfile, err := createTempFile(content)
		if err != nil {
			t.Fatal(err)
		}

		defer os.Remove(tmpfile.Name())

		superDatabaseUrl := "sqlite://super.db"
		crocoderDatabaseUrl := "sqlite://crocoder.db"
		acmeDatabaseUrl := "sqlite://acme.db"
		expected := &TomlConfig{
			SuperDatabaseUrl: &superDatabaseUrl,
			Tenants: map[string]Tenant{
				"crocoder": {Name: "crocoder", SubdomainName: "crocoder", DatabaseUrl: &crocoderDatabaseUrl},
				"acme":     {DatabaseUrl: &acmeDatabaseUrl},
			},
		}

		cfg, err := LoadConfigToml(tmpfile.Name())
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if !reflect.DeepEqual(cfg, expected) {
			t.Errorf("Expected %+v, got %+v", expected, cfg)
		}

	})

	t.Run("empty_toml", func(t *testing.T) {
		content := ""
		tmpfile, err := createTempFile(content)
		if err != nil {
			t.Fatal(err)
		}

		defer os.Remove(tmpfile.Name())

		expected := &TomlConfig{}

		cfg, err := LoadConfigToml(tmpfile.Name())
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if !reflect.DeepEqual(cfg, expected) {
			t.Errorf("Expected %+v, got %+v", expected, cfg)
		}
	})

	t.Run("tenant_database_template_example", func(t *testing.T) {
		content := `
super_database_url = "sqlite://super.db"
tenant_database_url_template = "sqlite://%s.db"
`
		tmpfile, err := createTempFile(content)
		if err != nil {
			t.Fatal(err)
		}

		defer os.Remove(tmpfile.Name())

		superDatabaseUrl := "sqlite://super.db"
		tenantDatabaseUrlTemplate := "sqlite://%s.db"
		expected := &TomlConfig{
			SuperDatabaseUrl:          &superDatabaseUrl,
			TenantDatabaseUrlTemplate: &tenantDatabaseUrlTemplate,
		}

		cfg, err := LoadConfigToml(tmpfile.Name())
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if !reflect.DeepEqual(cfg, expected) {
			t.Errorf("Expected %+v, got %+v", expected, cfg)
		}

	})

}

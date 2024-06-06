package config

import (
	"errors"
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   *config
		want    Config
		wantErr error
	}{
		{
			name: "SuperDatabaseUrl provided",
			input: &config{
				superDatabaseUrl: stringPtr("http://super.database.url"),
			},
			want: Config{
				IsMultiTenant:          true,
				IsSuperDatabaseEnabled: true,
				SuperDatabaseUrl:       "http://super.database.url",
			},
			wantErr: nil,
		},
		{
			name: "Valid tenants provided",
			input: &config{
				tenants: map[string]tenant{
					"tenant1": {
						subdomain:        stringPtr("tenant1"),
						databaseFilePath: stringPtr("/path/to/db1"),
					},
					"tenant2": {
						subdomain:   stringPtr("tenant2"),
						databaseUrl: stringPtr("http://tenant2.database.url"),
					},
				},
			},
			want: Config{
				IsMultiTenant: true,
				tenants: []Tenant{
					{
						Name:              "tenant1",
						Subdomain:         "tenant1",
						DatabaseType:      SQLite,
						LocalDatabasePath: "/path/to/db1",
					},
					{
						Name:         "tenant2",
						Subdomain:    "tenant2",
						DatabaseType: LibSQL,
						DatabaseUrl:  "http://tenant2.database.url",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:    "Both superDatabaseUrl and tenants empty",
			input:   &config{},
			want:    Config{},
			wantErr: errors.New("both super database url and tenants cannot be empty"),
		},
		{
			name: "Tenant with both databaseUrl and databaseFilePath empty",
			input: &config{
				tenants: map[string]tenant{
					"tenant1": {},
				},
			},
			want:    Config{},
			wantErr: errors.New("both database url and file path cannot be empty"),
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getConfig(tt.input)
			if (err != nil) && (tt.wantErr == nil || err.Error() != tt.wantErr.Error()) {
				t.Errorf("getConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

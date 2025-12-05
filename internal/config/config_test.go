package config_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tommarien/movie-land/internal/config"
)

func TestConfig(t *testing.T) {
	defaultEnvVars := map[string]string{
		"DATABASE_URL": "postgres://user:pass@localhost:5432/movie-land",
	}

	tests := []struct {
		name    string
		envVars map[string]string
		wantCfg *config.Config
		wantErr string
	}{
		{
			name: "return a config with the DATABASE_URL env var",
			envVars: map[string]string{
				"DATABASE_URL": "postgres://user:pass@localhost:5432/video-land",
			},
			wantCfg: &config.Config{
				DatabaseUrl: "postgres://user:pass@localhost:5432/video-land",
				Port:        3000,
			},
		},
		{
			name: "return a config with the PORT env var if set",
			envVars: map[string]string{
				"PORT": "4004",
			},
			wantCfg: &config.Config{
				DatabaseUrl: "postgres://user:pass@localhost:5432/movie-land",
				Port:        4004,
			},
		},
		{
			name: "returns an error when DATABASE_URL is not set",
			envVars: map[string]string{
				"DATABASE_URL": "",
			},
			wantErr: "environment variable \"DATABASE_URL\" should not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range defaultEnvVars {
				t.Setenv(k, v)
			}

			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			cfg, err := config.FromEnv()
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}

				if want := fmt.Sprintf("env: %s", tt.wantErr); err.Error() != want {
					t.Fatalf("expected error %q but got %q", want, err.Error())
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.wantCfg, cfg); diff != "" {
				t.Fatalf("config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

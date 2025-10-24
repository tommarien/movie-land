package config_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommarien/movie-land/internal/config"
)

func TestConfig(t *testing.T) {
	defaultEnvVars := map[string]string{
		"DATABASE_URL": "postgres://user:pass@localhost:5432/movie-land",
	}

	tests := []struct {
		name    string
		envVars map[string]string
		want    *config.Config
		wantErr string
	}{
		{
			name: "return a config",
			envVars: map[string]string{
				"DATABASE_URL": "postgres://user:pass@localhost:5432/video-land",
			},
			want: &config.Config{
				DatabaseUrl: "postgres://user:pass@localhost:5432/video-land",
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
				assert.EqualError(t, err, fmt.Sprintf("env: %s", tt.wantErr))
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, cfg)
		})
	}
}

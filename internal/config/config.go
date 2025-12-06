package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseUrl         string        `env:"DATABASE_URL,notEmpty"`
	DatabasePingTimeout time.Duration `env:"DATABASE_PING_TIMEOUT" envDefault:"200ms"`
	Port                int           `env:"PORT" envDefault:"3000"`
}

func FromEnv() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

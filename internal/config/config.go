package config

import "github.com/caarlos0/env/v11"

type Config struct {
	DatabaseUrl string `env:"DATABASE_URL,notEmpty"`
}

func FromEnv() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

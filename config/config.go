package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/vigorouzis/aibolit-notification/internal/infrastructure/postgres"
	"github.com/vigorouzis/aibolit-notification/internal/interface/http"
)

type Config struct {
	Postgres postgres.Config `envPrefix:"POSTGRES_"`
	HTTP     http.Config     `envPrefix:"HTTP_"`
}

func FromENV() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}
	
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

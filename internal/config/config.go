package config

import (
	"log"
	"socket/internal/database"
	"socket/internal/server"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Server server.Config   `envPrefix:"SERVER_"`
	DB     database.Config `envPrefix:"DB_"`
}

// Load configs from .env file
func Load() *AppConfig {
	config := &AppConfig{}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Unable to load .env file: %s", err)
	}

	if err := env.Parse(config); err != nil {
		log.Fatalf("Unable to parse env vars: %s", err)
	}

	return config
}

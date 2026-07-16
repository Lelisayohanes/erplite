package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Host string
	Port string
}

func Load() (*Config, error) {
	// Load .env if it exists.
	// Ignore the error because production environments may not use a .env file.
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST"),
			Port: getEnv("SERVER_PORT"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("SERVER_HOST environment variable is required")
	}
	if c.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT environment variable is required")
	}
	return nil
}

func getEnv(key string) string {
	return os.Getenv(key)
}

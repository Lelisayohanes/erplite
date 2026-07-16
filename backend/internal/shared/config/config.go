package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTP HTTPConfig
}

type HTTPConfig struct {
	Host string
	Port string

	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration

	MaxHeaderBytes int
}

func Load() (*Config, error) {

	// Load .env during development.
	// Production environments normally provide real environment variables.
	_ = godotenv.Load()

	cfg := &Config{
		HTTP: HTTPConfig{
			Host: getEnv("HTTP_HOST"),
			Port: getEnv("HTTP_PORT"),

			ReadTimeout: getDuration(
				"HTTP_READ_TIMEOUT",
				15*time.Second,
			),

			WriteTimeout: getDuration(
				"HTTP_WRITE_TIMEOUT",
				15*time.Second,
			),

			IdleTimeout: getDuration(
				"HTTP_IDLE_TIMEOUT",
				60*time.Second,
			),

			ReadHeaderTimeout: getDuration(
				"HTTP_READ_HEADER_TIMEOUT",
				5*time.Second,
			),

			ShutdownTimeout: getDuration(
				"HTTP_SHUTDOWN_TIMEOUT",
				10*time.Second,
			),

			MaxHeaderBytes: getInt(
				"HTTP_MAX_HEADER_BYTES",
				1048576,
			),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {

	if c.HTTP.Host == "" {
		return fmt.Errorf(
			"HTTP_HOST environment variable is required",
		)
	}

	if c.HTTP.Port == "" {
		return fmt.Errorf(
			"HTTP_PORT environment variable is required",
		)
	}

	if c.HTTP.ShutdownTimeout <= 0 {
		return fmt.Errorf(
			"shutdown timeout must be greater than zero",
		)
	}

	if c.HTTP.MaxHeaderBytes <= 0 {
		return fmt.Errorf(
			"max header bytes must be greater than zero",
		)
	}

	return nil
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func getDuration(key string, fallback time.Duration) time.Duration {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)

	if err != nil {
		return fallback
	}

	return duration
}

func getInt(key string, fallback int) int {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	result, err := strconv.Atoi(value)

	if err != nil {
		return fallback
	}

	return result
}

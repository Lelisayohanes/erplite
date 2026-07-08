package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all application configuration.
// All required values must be provided via the .env file or environment
// variables.  There are no built-in defaults — the application will
// refuse to start if any required value is missing.
//
// Supported sources (highest → lowest precedence):
//
//  1. Environment variables (prefixed with ERPLITE_, e.g. ERPLITE_DATABASE_HOST)
//  2. .env file (loaded by godotenv; does NOT overwrite existing env vars)
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	// add others: Redis, SMTP, etc.
}

type AppConfig struct {
	Name    string
	Env     string
	Port    string
	Debug   bool
	Timeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Load reads configuration from supported external sources and applies
// the defined precedence order.  There are no built-in defaults — every
// required value must be supplied via the .env file or environment
// variables.  Returns an error if any required value is missing.
func Load() (*Config, error) {
	// ── Layer 2: .env file ───────────────────────────────────────────
	// Loads key=value pairs into the process environment.
	// Existing env vars are NOT overwritten by godotenv.
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("config: failed to load .env file: %w", err)
	}

	// ── Layer 1: Environment variables (highest precedence) ──────────
	viper.SetEnvPrefix("ERPLITE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Bind every config key explicitly so viper reads the corresponding
	// environment variable (e.g. ERPLITE_APP_PORT → app.port).
	for _, key := range allKeys {
		if err := viper.BindEnv(key); err != nil {
			return nil, fmt.Errorf("config: failed to bind key %q: %w", key, err)
		}
	}

	// Unmarshal into struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: failed to unmarshal: %w", err)
	}

	// Parse timeout duration from the resolved string value
	if d, err := time.ParseDuration(viper.GetString("app.timeout")); err == nil {
		cfg.App.Timeout = d
	} else {
		return nil, fmt.Errorf("config: invalid app.timeout value: %w", err)
	}

	// Validate required fields — fails if any value is missing
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// allKeys lists every configuration key the application requires.
var allKeys = []string{
	"app.name",
	"app.env",
	"app.port",
	"app.debug",
	"app.timeout",
	"database.host",
	"database.port",
	"database.user",
	"database.password",
	"database.dbname",
	"database.sslmode",
}

// Validate checks that all required configuration values are present and
// returns a descriptive error listing every missing field.
func (c *Config) Validate() error {
	var missing []string

	if c.App.Name == "" {
		missing = append(missing, "app.name (ERPLITE_APP_NAME)")
	}
	if c.App.Env == "" {
		missing = append(missing, "app.env (ERPLITE_APP_ENV)")
	}
	if c.App.Port == "" {
		missing = append(missing, "app.port (ERPLITE_APP_PORT)")
	}
	if c.Database.Host == "" {
		missing = append(missing, "database.host (ERPLITE_DATABASE_HOST)")
	}
	if c.Database.Port == "" {
		missing = append(missing, "database.port (ERPLITE_DATABASE_PORT)")
	}
	if c.Database.User == "" {
		missing = append(missing, "database.user (ERPLITE_DATABASE_USER)")
	}
	if c.Database.DBName == "" {
		missing = append(missing, "database.dbname (ERPLITE_DATABASE_DBNAME)")
	}
	if c.Database.SSLMode == "" {
		missing = append(missing, "database.sslmode (ERPLITE_DATABASE_SSLMODE)")
	}

	if len(missing) > 0 {
		return fmt.Errorf(
			"config: required configuration values are missing:\n  - %s",
			strings.Join(missing, "\n  - "),
		)
	}
	return nil
}

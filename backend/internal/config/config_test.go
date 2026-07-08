package config

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// setAllEnv sets every required config env var with the given values.
func setAllEnv(t *testing.T) {
	t.Helper()
	envs := map[string]string{
		"ERPLITE_APP_NAME":          "test-app",
		"ERPLITE_APP_ENV":           "test",
		"ERPLITE_APP_PORT":          "9090",
		"ERPLITE_APP_DEBUG":         "true",
		"ERPLITE_APP_TIMEOUT":       "10s",
		"ERPLITE_APP_LOGFORMAT":     "json",
		"ERPLITE_DATABASE_HOST":     "localhost",
		"ERPLITE_DATABASE_PORT":     "5432",
		"ERPLITE_DATABASE_USER":     "testuser",
		"ERPLITE_DATABASE_PASSWORD": "testpass",
		"ERPLITE_DATABASE_DBNAME":   "testdb",
		"ERPLITE_DATABASE_SSLMODE":  "disable",
	}
	for k, v := range envs {
		t.Setenv(k, v)
	}
}

func resetViper(t *testing.T) {
	t.Helper()
	viper.Reset()
}

// ensureDotEnv creates a temporary empty .env if one doesn't exist,
// so that godotenv.Load() succeeds in CI where no .env is present.
// It restores the original after the test.
func ensureDotEnv(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(".env"); err != nil {
		os.WriteFile(".env", []byte{}, 0o644)
		t.Cleanup(func() { os.Remove(".env") })
	}
}

func TestLoad_Success(t *testing.T) {
	resetViper(t)
	setAllEnv(t)
	ensureDotEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.App.Name != "test-app" {
		t.Errorf("expected App.Name=test-app, got %s", cfg.App.Name)
	}
	if cfg.App.Port != "9090" {
		t.Errorf("expected App.Port=9090, got %s", cfg.App.Port)
	}
	if cfg.App.Timeout != 10*time.Second {
		t.Errorf("expected App.Timeout=10s, got %v", cfg.App.Timeout)
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("expected Database.Host=localhost, got %s", cfg.Database.Host)
	}
}

func TestLoad_MissingRequired(t *testing.T) {
	resetViper(t)
	// Clear all required env vars
	for _, key := range []string{
		"ERPLITE_APP_NAME", "ERPLITE_APP_ENV", "ERPLITE_APP_PORT",
		"ERPLITE_APP_LOGFORMAT",
		"ERPLITE_DATABASE_HOST", "ERPLITE_DATABASE_PORT",
		"ERPLITE_DATABASE_USER", "ERPLITE_DATABASE_PASSWORD",
		"ERPLITE_DATABASE_DBNAME", "ERPLITE_DATABASE_SSLMODE",
	} {
		os.Unsetenv(key)
	}
	ensureDotEnv(t)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
	if !strings.Contains(err.Error(), "config:") {
		t.Errorf("expected error to mention 'config:', got: %v", err)
	}
}

func TestLoad_EnvOverridesDotEnv(t *testing.T) {
	resetViper(t)
	setAllEnv(t)
	ensureDotEnv(t)
	// Override with env var — should win
	t.Setenv("ERPLITE_APP_PORT", "7777")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.App.Port != "7777" {
		t.Errorf("expected env var override port=7777, got %s", cfg.App.Port)
	}
}

func TestDSN(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     "db.example.com",
			Port:     "5433",
			User:     "admin",
			Password: "secret",
			DBName:   "mydb",
			SSLMode:  "require",
		},
	}
	want := "postgres://admin:secret@db.example.com:5433/mydb?sslmode=require"
	got := cfg.DSN()
	if got != want {
		t.Errorf("DSN() = %q, want %q", got, want)
	}
}

func TestValidate_AllPresent(t *testing.T) {
	cfg := &Config{
		App: AppConfig{
			Name: "app", Env: "test", Port: "8080", LogFormat: "json",
		},
		Database: DatabaseConfig{
			Host: "h", Port: "5432", User: "u", Password: "p",
			DBName: "db", SSLMode: "disable",
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingFields(t *testing.T) {
	cfg := &Config{} // all zero values
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	msg := err.Error()
	for _, field := range []string{
		"app.name", "app.env", "app.port", "app.logformat",
		"database.host", "database.port", "database.user",
		"database.dbname", "database.sslmode",
	} {
		if !strings.Contains(msg, field) {
			t.Errorf("expected error to mention %q", field)
		}
	}
}

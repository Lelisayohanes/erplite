package features

import (
	"os"
	"testing"

	"github.com/cucumber/godog"

	"erplite/backend/features/steps"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: steps.InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"."},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func TestMain(m *testing.M) {
	// Set required env vars for BDD steps that call config.Load().
	// In CI these would come from GitHub Actions secrets; here we provide
	// safe test defaults so scenarios can exercise the real config loader.
	// Viper reads the .env file natively (optional), and real env vars
	// always take precedence.
	testEnv := map[string]string{
		"ERPLITE_APP_NAME":          "erplite",
		"ERPLITE_APP_ENV":           "test",
		"ERPLITE_APP_PORT":          "8080",
		"ERPLITE_APP_DEBUG":         "true",
		"ERPLITE_APP_TIMEOUT":       "30s",
		"ERPLITE_APP_LOGFORMAT":     "text",
		"ERPLITE_DATABASE_HOST":     "localhost",
		"ERPLITE_DATABASE_PORT":     "5432",
		"ERPLITE_DATABASE_USER":     "postgres",
		"ERPLITE_DATABASE_PASSWORD": "postgres",
		"ERPLITE_DATABASE_DBNAME":   "erplite",
		"ERPLITE_DATABASE_SSLMODE":  "disable",
	}
	for k, v := range testEnv {
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}

	code := m.Run()

	os.Exit(code)
}

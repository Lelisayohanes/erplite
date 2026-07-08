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
	// Ensure a .env file exists for config.Load() calls in BDD steps.
	// In CI there is no .env file; godotenv.Load() would fail without one.
	// We create a temporary empty .env in the backend/ directory (parent of
	// features/) and clean it up after all tests finish.
	envPath := "../.env"
	created := false
	if _, err := os.Stat(envPath); err != nil {
		if f, err := os.Create(envPath); err == nil {
			f.Close()
			created = true
		}
	}

	code := m.Run()

	if created {
		os.Remove(envPath)
	}
	os.Exit(code)
}

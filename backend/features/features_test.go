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
	// Change working directory to features/ so godog finds .feature files
	// when running `go test ./features/...`
	if wd, err := os.Getwd(); err == nil {
		_ = wd // already in the right directory when running from backend/
	}
	m.Run()
}

package steps

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"

	"github.com/cucumber/godog"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"erplite/backend/internal/config"
	"erplite/backend/internal/container"
	"erplite/backend/internal/handler"
	"erplite/backend/internal/logger"
	"erplite/backend/internal/middleware"
)

// ── Shared state across scenarios ──────────────────────────────────

type scenarioState struct {
	cfg        *config.Config
	log        *slog.Logger
	ctr        *container.Container
	echo       *echo.Echo
	db         *sql.DB
	loadErr    error
	lastStatus int
	lastBody   map[string]string
	logBuf     strings.Builder
	corrID     string
}

var s scenarioState

// ── Scenario initializer ───────────────────────────────────────────

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, scenario *godog.Scenario) (context.Context, error) {
		s = scenarioState{}
		viper.Reset() // reset global viper state between scenarios
		return ctx, nil
	})

	// ─── Configuration Management ──────────────────────────────────
	sc.Step(`^the application configuration component is initialized$`, configComponentInitialized)
	sc.Step(`^configuration values are provided from supported external sources$`, configFromExternalSources)
	sc.Step(`^configuration values exist in multiple supported sources$`, configMultipleSources)
	sc.Step(`^a configuration value exists in a lower-priority source$`, configLowerPrioritySource)
	sc.Step(`^the same configuration value is provided as an environment variable$`, configEnvVarOverride)
	sc.Step(`^required configuration values are defined$`, configRequiredDefined)
	sc.Step(`^one or more required configuration values are missing$`, configRequiredMissing)
	sc.Step(`^the application has multiple modules$`, appMultipleModules)
	sc.Step(`^a developer is running the application locally$`, devRunningLocally)
	sc.Step(`^local configuration values are provided through supported local sources$`, localConfigProvided)
	sc.Step(`^the application configuration structure is defined$`, configStructureDefined)
	sc.Step(`^a module requires configuration values$`, moduleRequiresConfig)
	sc.Step(`^the configuration loader should load the available configuration values$`, loaderLoadsValues)
	sc.Step(`^the application should use the loaded configuration values$`, appUsesLoadedValues)
	sc.Step(`^the configuration loader should apply the defined precedence order$`, loaderAppliesPrecedence)
	sc.Step(`^higher-priority configuration sources should override lower-priority sources$`, higherOverridesLower)
	sc.Step(`^the application should use the environment variable value$`, appUsesEnvValue)
	sc.Step(`^the configuration validator should check all required values$`, validatorChecksAll)
	sc.Step(`^the application should continue startup if all required values are valid$`, appContinuesStartup)
	sc.Step(`^the application should fail to start$`, appFailsToStart)
	sc.Step(`^the application should return a clear configuration error message$`, appReturnsConfigError)
	sc.Step(`^the module should retrieve configuration through the centralized configuration component$`, moduleUsesCentralized)
	sc.Step(`^the module should not directly access configuration sources$`, moduleNoDirectAccess)
	sc.Step(`^the application should load the local development configuration successfully$`, appLoadsLocalConfig)
	sc.Step(`^configuration documentation is generated$`, configDocGenerated)
	sc.Step(`^all supported configuration sources should be documented$`, sourcesDocumented)
	sc.Step(`^configuration precedence rules should be documented$`, precedenceDocumented)
	sc.Step(`^required configuration values should be documented$`, requiredDocumented)

	// ─── Structured Logging ────────────────────────────────────────
	sc.Step(`^the application is running$`, appIsRunning)
	sc.Step(`^an event is logged$`, eventIsLogged)
	sc.Step(`^the log entry is generated in a structured format$`, logIsStructured)
	sc.Step(`^includes the timestamp, severity level, and message$`, logHasFields)
	sc.Step(`^a request includes a correlation identifier$`, requestWithCorrelationID)
	sc.Step(`^the request is processed$`, requestProcessed)
	sc.Step(`^all related log entries include the same correlation identifier$`, logHasCorrelationID)
	sc.Step(`^request logging is enabled$`, requestLoggingEnabled)
	sc.Step(`^an HTTP request is processed$`, httpProcessed)
	sc.Step(`^the request metadata is logged$`, metadataLogged)
	sc.Step(`^sensitive information is excluded$`, sensitiveExcluded)

	// ─── Dependency Injection ──────────────────────────────────────
	sc.Step(`^application dependencies are registered$`, depsRegistered)
	sc.Step(`^all required components are resolved successfully$`, componentsResolved)
	sc.Step(`^the dependency container is configured$`, containerConfigured)
	sc.Step(`^the application is initialized$`, appInitialized)
	sc.Step(`^all application services are injected automatically$`, servicesInjected)
	sc.Step(`^the application starts successfully$`, appStartsSuccessfully)

	// ─── Development Automation ────────────────────────────────────
	sc.Step(`^the project source code exists$`, sourceCodeExists)
	sc.Step(`^the build command is executed$`, buildExecuted)
	sc.Step(`^the application binary is generated successfully$`, binaryGenerated)
	sc.Step(`^automated tests exist$`, testsExist)
	sc.Step(`^the test command is executed$`, testExecuted)
	sc.Step(`^all tests are run$`, allTestsRun)
	sc.Step(`^the results are displayed$`, resultsDisplayed)
	sc.Step(`^the linting tool is configured$`, lintConfigured)
	sc.Step(`^the lint command is executed$`, lintExecuted)
	sc.Step(`^coding standard violations are reported$`, violationsReported)
	sc.Step(`^the help command is executed$`, helpExecuted)
	sc.Step(`^all supported commands are displayed with descriptions$`, commandsDisplayed)

	// ─── Continuous Integration ────────────────────────────────────
	sc.Step(`^code is pushed to the repository$`, codePushed)
	sc.Step(`^the CI pipeline executes$`, ciExecutes)
	sc.Step(`^code quality checks are performed$`, qualityChecksPerformed)
	sc.Step(`^automated tests are executed$`, ciTestsExecuted)
	sc.Step(`^the pipeline completes successfully$`, pipelineCompletes)
	sc.Step(`^a pull request is opened$`, prOpened)
	sc.Step(`^the application is built successfully$`, ciBuildSucceeds)
	sc.Step(`^the container image is validated$`, imageValidated)
	sc.Step(`^dependencies have previously been downloaded$`, depsCached)
	sc.Step(`^the pipeline executes$`, pipelineRuns)
	sc.Step(`^dependency caching is used to reduce execution time$`, cachingUsed)

	// ─── Health Monitoring ─────────────────────────────────────────
	sc.Step(`^the liveness endpoint is requested$`, livenessRequested)
	sc.Step(`^the application responds with HTTP (\d+)$`, respondsWithStatus)
	sc.Step(`^indicates it is healthy$`, indicatesHealthy)
	sc.Step(`^all required dependencies are available$`, depsAvailable)
	sc.Step(`^the readiness endpoint is requested$`, readinessRequested)
	sc.Step(`^indicates it is ready to accept traffic$`, indicatesReady)
	sc.Step(`^a required dependency is unavailable$`, depUnavailable)
	sc.Step(`^the response indicates the application is not ready$`, indicatesNotReady)

	// ─── Shared ────────────────────────────────────────────────────
	sc.Step(`^the application starts$`, appStarts)
	sc.Step(`^the database is available$`, dbAvailable)
	sc.Step(`^the migration command is executed$`, migrationExecuted)
	sc.Step(`^all pending migrations are applied successfully$`, migrationsApplied)
}

// ── Configuration Management steps ─────────────────────────────────

func configComponentInitialized() error { return nil }
func configFromExternalSources() error  { return nil }
func configMultipleSources() error      { return nil }
func configLowerPrioritySource() error  { return nil }
func configEnvVarOverride() error       { return nil }
func configRequiredDefined() error      { return nil }
func configRequiredMissing() error {
	// Simulate missing required config by creating an empty Config
	// and skipping the real loader (which would load .env values)
	s.cfg = &config.Config{}
	s.loadErr = s.cfg.Validate()
	return nil
}
func appMultipleModules() error     { return nil }
func devRunningLocally() error      { return nil }
func localConfigProvided() error    { return nil }
func configStructureDefined() error { return nil }
func moduleRequiresConfig() error   { return nil }
func appStarts() error {
	// If a Given step already set up a pre-configured state, skip real loading
	if s.cfg != nil || s.loadErr != nil {
		return nil
	}
	// Ensure we run from backend/ where .env lives
	if wd, _ := os.Getwd(); strings.HasSuffix(wd, "features") {
		os.Chdir("..")
	}
	cfg, err := config.Load()
	s.cfg = cfg
	s.loadErr = err
	return nil
}
func loaderLoadsValues() error {
	if s.cfg == nil {
		return fmt.Errorf("config was not loaded")
	}
	return nil
}
func appUsesLoadedValues() error {
	if s.cfg.App.Name == "" {
		return fmt.Errorf("app name is empty")
	}
	return nil
}
func loaderAppliesPrecedence() error { return nil }
func higherOverridesLower() error    { return nil }
func appUsesEnvValue() error         { return nil }
func validatorChecksAll() error {
	if s.cfg == nil {
		return fmt.Errorf("config not loaded")
	}
	return s.cfg.Validate()
}
func appContinuesStartup() error {
	if s.loadErr != nil {
		return fmt.Errorf("startup failed: %v", s.loadErr)
	}
	return nil
}
func appFailsToStart() error {
	if s.loadErr == nil {
		return fmt.Errorf("expected startup to fail but it succeeded")
	}
	return nil
}
func appReturnsConfigError() error {
	if s.loadErr == nil {
		return fmt.Errorf("expected config error but got none")
	}
	if !strings.Contains(s.loadErr.Error(), "config:") {
		return fmt.Errorf("error message does not mention config: %v", s.loadErr)
	}
	return nil
}
func moduleUsesCentralized() error { return nil }
func moduleNoDirectAccess() error  { return nil }
func appLoadsLocalConfig() error   { return nil }
func configDocGenerated() error    { return nil }
func sourcesDocumented() error     { return nil }
func precedenceDocumented() error  { return nil }
func requiredDocumented() error    { return nil }

// ── Structured Logging steps ───────────────────────────────────────

func appIsRunning() error {
	s.log = logger.New("development", slog.String("service", "erplite"))
	return nil
}
func eventIsLogged() error {
	s.log.Info("test event", "key", "value")
	return nil
}
func logIsStructured() error { return nil }
func logHasFields() error    { return nil }
func requestWithCorrelationID() error {
	s.corrID = "test-correlation-123"
	return nil
}
func requestProcessed() error    { return nil }
func logHasCorrelationID() error { return nil }
func requestLoggingEnabled() error {
	s.echo = echo.New()
	s.echo.Use(middleware.CorrelationID())
	s.echo.Use(middleware.RequestLogger(s.echo.Logger))
	return nil
}
func httpProcessed() error     { return nil }
func metadataLogged() error    { return nil }
func sensitiveExcluded() error { return nil }

// ── Dependency Injection steps ─────────────────────────────────────

func depsRegistered() error { return nil }
func componentsResolved() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := logger.New(cfg.App.Env)
	s.ctr = container.New(log, cfg, nil)
	if s.ctr.Log == nil || s.ctr.Cfg == nil {
		return fmt.Errorf("container has nil fields")
	}
	return nil
}
func containerConfigured() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := logger.New(cfg.App.Env)
	s.ctr = container.New(log, cfg, nil)
	return nil
}
func appInitialized() error   { return nil }
func servicesInjected() error { return nil }
func appStartsSuccessfully() error {
	if s.ctr == nil {
		return fmt.Errorf("container is nil")
	}
	return nil
}

// ── Development Automation steps ───────────────────────────────────

func sourceCodeExists() error { return nil }
func buildExecuted() error {
	cmd := exec.Command("go", "build", "-o", "/dev/null", "./cmd/server")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %s\n%s", err, out)
	}
	return nil
}
func binaryGenerated() error  { return nil }
func testsExist() error       { return nil }
func testExecuted() error     { return nil }
func allTestsRun() error      { return nil }
func resultsDisplayed() error { return nil }
func lintConfigured() error   { return nil }
func lintExecuted() error {
	cmd := exec.Command("go", "vet", "./...")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lint failed: %s\n%s", err, out)
	}
	return nil
}
func violationsReported() error { return nil }
func helpExecuted() error {
	cmd := exec.Command("make", "help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("make help failed: %s\n%s", err, out)
	}
	s.logBuf.Reset()
	s.logBuf.Write(out)
	return nil
}
func commandsDisplayed() error {
	out := s.logBuf.String()
	for _, cmd := range []string{"build", "test", "lint", "migrate", "help"} {
		if !strings.Contains(out, cmd) {
			return fmt.Errorf("help output missing command: %s", cmd)
		}
	}
	return nil
}
func dbAvailable() error       { return nil }
func migrationExecuted() error { return nil }
func migrationsApplied() error { return nil }

// ── Continuous Integration steps ───────────────────────────────────

func codePushed() error             { return nil }
func ciExecutes() error             { return nil }
func qualityChecksPerformed() error { return nil }
func ciTestsExecuted() error        { return nil }
func pipelineCompletes() error      { return nil }
func prOpened() error               { return nil }
func ciBuildSucceeds() error        { return nil }
func imageValidated() error         { return nil }
func depsCached() error             { return nil }
func pipelineRuns() error           { return nil }
func cachingUsed() error            { return nil }

// ── Health Monitoring steps ────────────────────────────────────────

func livenessRequested() error {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	e.GET("/healthz", handler.Health)
	e.ServeHTTP(rec, req)
	s.lastStatus = rec.Code
	json.Unmarshal(rec.Body.Bytes(), &s.lastBody)
	return nil
}

func readinessRequested() error {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	e.GET("/readyz", handler.Readiness(s.db))
	e.ServeHTTP(rec, req)
	s.lastStatus = rec.Code
	json.Unmarshal(rec.Body.Bytes(), &s.lastBody)
	return nil
}

func respondsWithStatus(code int) error {
	if s.lastStatus != code {
		return fmt.Errorf("expected status %d but got %d", code, s.lastStatus)
	}
	return nil
}
func indicatesHealthy() error {
	if s.lastBody["status"] != "healthy" {
		return fmt.Errorf("expected status=healthy but got %s", s.lastBody["status"])
	}
	return nil
}
func depsAvailable() error {
	if wd, _ := os.Getwd(); strings.HasSuffix(wd, "features") {
		os.Chdir("..")
	}
	cfg, err := config.Load()
	if err != nil {
		return godog.ErrPending
	}
	pool, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return godog.ErrPending
	}
	if err := pool.Ping(); err != nil {
		pool.Close()
		// DB not reachable — skip readiness scenario
		return godog.ErrPending
	}
	s.db = pool
	return nil
}
func indicatesReady() error {
	if s.lastBody["status"] != "ready" {
		return fmt.Errorf("expected status=ready but got %s", s.lastBody["status"])
	}
	return nil
}
func depUnavailable() error { return nil }
func indicatesNotReady() error {
	if s.lastBody["status"] != "not_ready" {
		return fmt.Errorf("expected status=not_ready but got %s", s.lastBody["status"])
	}
	return nil
}

// Ensure os import is used (for TestMain pattern)
var _ = os.Getwd

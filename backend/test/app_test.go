package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lelisayohanes/erplite/backend/internal/app"
	"github.com/Lelisayohanes/erplite/backend/internal/shared/config"
)

func testConfig() *config.Config {
	return &config.Config{
		HTTP: config.HTTPConfig{
			Host: "localhost",
			Port: "8080",
		},
	}
}

func createTestApp() *app.App {
	return app.New(testConfig())
}

func performRequest(
	application *app.App,
	method string,
	path string,
) *httptest.ResponseRecorder {

	req := httptest.NewRequest(
		method,
		path,
		nil,
	)

	recorder := httptest.NewRecorder()

	application.Handler().ServeHTTP(
		recorder,
		req,
	)

	return recorder
}

func TestAppCreation(t *testing.T) {

	application := createTestApp()

	if application == nil {
		t.Fatal("expected application, got nil")
	}
}

func TestHealthEndpoint(t *testing.T) {

	response := performRequest(
		createTestApp(),
		http.MethodGet,
		"/health",
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}
}

func TestPingEndpoint(t *testing.T) {

	response := performRequest(
		createTestApp(),
		http.MethodGet,
		"/ping",
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}
}

func TestReadyEndpoint(t *testing.T) {

	response := performRequest(
		createTestApp(),
		http.MethodGet,
		"/ready",
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}
}

func TestLiveEndpoint(t *testing.T) {

	response := performRequest(
		createTestApp(),
		http.MethodGet,
		"/live",
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}
}

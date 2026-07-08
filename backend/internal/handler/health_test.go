package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHealth(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	if err := Health(c); err != nil {
		t.Fatalf("Health() returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if body == "" {
		t.Error("expected non-empty body")
	}
}

func TestReadiness_NilDB(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	h := Readiness(nil)
	if err := h(c); err != nil {
		t.Fatalf("Readiness(nil) returned error: %v", err)
	}
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", rec.Code)
	}
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestCorrelationID_GeneratesWhenMissing(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e.Use(CorrelationID())
	e.GET("/", func(c echo.Context) error {
		id, _ := c.Get("correlation_id").(string)
		if id == "" {
			t.Error("expected correlation_id to be set in context")
		}
		return c.String(http.StatusOK, id)
	})

	e.ServeHTTP(rec, req)

	if rec.Header().Get(echo.HeaderXCorrelationID) == "" {
		t.Error("expected X-Correlation-ID response header to be set")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCorrelationID_PreservesIncoming(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderXCorrelationID, "incoming-id-42")
	rec := httptest.NewRecorder()

	e.Use(CorrelationID())
	e.GET("/", func(c echo.Context) error {
		id, _ := c.Get("correlation_id").(string)
		if id != "incoming-id-42" {
			t.Errorf("expected correlation_id=incoming-id-42, got %s", id)
		}
		return c.String(http.StatusOK, "ok")
	})

	e.ServeHTTP(rec, req)

	if rec.Header().Get(echo.HeaderXCorrelationID) != "incoming-id-42" {
		t.Error("expected response header to echo incoming correlation ID")
	}
}

func TestRequestLogger_DoesNotPanic(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	e.Use(CorrelationID())
	e.Use(RequestLogger(e.Logger))
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// Should not panic
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

package logger

import (
	"log/slog"
	"testing"
)

func TestNew_Development(t *testing.T) {
	log := New("development")
	if log == nil {
		t.Fatal("expected non-nil logger")
	}
	// Development uses text handler — just verify it works
	log.Info("test message")
}

func TestNew_Production(t *testing.T) {
	log := New("production")
	if log == nil {
		t.Fatal("expected non-nil logger")
	}
	log.Info("test message")
}

func TestNew_WithAttributes(t *testing.T) {
	log := New("test",
		slog.String("service", "erplite"),
		slog.String("env", "test"),
	)
	if log == nil {
		t.Fatal("expected non-nil logger")
	}
	log.Info("test message with attributes")
}

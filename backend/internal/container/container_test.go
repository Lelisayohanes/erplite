package container

import (
	"log/slog"
	"testing"

	"erplite/backend/internal/config"
)

func TestNew(t *testing.T) {
	log := slog.Default()
	cfg := &config.Config{}

	ctr := New(log, cfg, nil)
	if ctr == nil {
		t.Fatal("expected non-nil container")
	}
	if ctr.Log == nil {
		t.Error("expected Log to be set")
	}
	if ctr.Cfg == nil {
		t.Error("expected Cfg to be set")
	}
	if ctr.DB != nil {
		t.Error("expected DB to be nil when not provided")
	}
}

func TestClose_NilDB(t *testing.T) {
	ctr := &Container{DB: nil}
	// Should not panic
	ctr.Close()
}

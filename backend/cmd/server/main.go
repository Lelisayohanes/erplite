package main

import (
	"log/slog"
	"os"
	"strconv"

	"erplite/backend/internal/config"
	"erplite/backend/internal/container"
	"erplite/backend/internal/db"
	"erplite/backend/internal/logger"
	"erplite/backend/internal/server"
)

func main() {
	// ── Load configuration ──────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		// No logger yet — fall back to stderr for config errors.
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	// ── Initialise structured logger ────────────────────────────────
	log := logger.New(cfg.App.Env,
		slog.String("service", cfg.App.Name),
		slog.String("env", cfg.App.Env),
	)

	// ── Connect to database ─────────────────────────────────────────
	pool, err := db.Connect(cfg.DSN(), log)
	if err != nil {
		log.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// ── Build dependency container ──────────────────────────────────
	ctr := container.New(log, cfg, pool)
	defer ctr.Close()

	// ── Handle CLI sub-commands ─────────────────────────────────────
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		dir := "up"
		steps := 0
		if len(os.Args) > 2 {
			dir = os.Args[2]
		}
		if len(os.Args) > 3 {
			steps, _ = strconv.Atoi(os.Args[3])
		}
		if err := db.RunMigrations(cfg.DSN(), dir, steps, log); err != nil {
			log.Error("migration failed", "error", err)
			os.Exit(1)
		}
		return
	}

	// ── Start server ────────────────────────────────────────────────
	log.Info("starting server", "port", cfg.App.Port)
	s := server.New(ctr)
	if err := s.Start(); err != nil {
		log.Error("server error", "error", err)
		os.Exit(1)
	}
}

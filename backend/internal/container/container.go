package container

import (
	"database/sql"
	"log/slog"

	"erplite/backend/internal/config"
)

// Container is a lightweight dependency container. All application
// dependencies are registered at startup and resolved from here.
type Container struct {
	Log *slog.Logger
	Cfg *config.Config
	DB  *sql.DB
}

// New creates a fully-resolved container. Every dependency is
// injected here so that downstream packages receive them via the
// container rather than constructing their own.
func New(log *slog.Logger, cfg *config.Config, db *sql.DB) *Container {
	return &Container{
		Log: log,
		Cfg: cfg,
		DB:  db,
	}
}

// Close releases resources held by managed dependencies.
func (c *Container) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}

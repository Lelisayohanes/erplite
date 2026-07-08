package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

// Connect opens a connection pool to PostgreSQL and verifies connectivity.
func Connect(dsn string, log *slog.Logger) (*sql.DB, error) {
	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("db: open connection: %w", err)
	}

	pool.SetMaxOpenConns(25)
	pool.SetMaxIdleConns(5)

	if err := pool.Ping(); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping failed: %w", err)
	}

	log.Info("database connected")
	return pool, nil
}

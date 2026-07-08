package db

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dsn string, direction string, steps int, log *slog.Logger) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("migration init failed: %w", err)
	}
	defer m.Close()

	switch direction {
	case "up":
		if steps == 0 {
			err = m.Up()
		} else {
			err = m.Steps(steps)
		}
	case "down":
		if steps == 0 {
			err = m.Down()
		} else {
			err = m.Steps(-steps)
		}
	default:
		return errors.New("direction must be 'up' or 'down'")
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	log.Info("migrations applied successfully")
	return nil
}

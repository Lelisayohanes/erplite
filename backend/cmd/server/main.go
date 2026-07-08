package main // Entry point of the Go application.

import (
	// Provides HTTP status codes like http.StatusOK.

	// Echo web framework.
	// Built-in middleware (logging, recovery, etc.).
	"erplite/backend/internal/db"
	"fmt"
	"log" // Logging package.
	"os"
	"strconv"

	"erplite/backend/internal/config" // Custom config package.
	"erplite/backend/internal/server" // Custom server package.
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.DBName, cfg.Database.SSLMode)

	// If the first argument is "migrate", run migration and exit
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		dir := "up"
		steps := 0
		if len(os.Args) > 2 {
			dir = os.Args[2]
		}
		if len(os.Args) > 3 {
			steps, _ = strconv.Atoi(os.Args[3])
		}
		if err := db.RunMigrations(dsn, dir, steps); err != nil {
			log.Fatal(err)
		}
		return
	}

	// otherwise start the server
	s := server.New(cfg)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}

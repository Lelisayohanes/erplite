package main

import (
	"log"

	"github.com/Lelisayohanes/erplite/backend/internal/app"
	"github.com/Lelisayohanes/erplite/backend/internal/shared/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	application := app.New(cfg)

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}

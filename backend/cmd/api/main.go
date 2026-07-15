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

	addr := cfg.Server.Host + ":" + cfg.Server.Port

	if err := application.Echo.Start(addr); err != nil {
		log.Fatal(err)
	}
}

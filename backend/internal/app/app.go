package app

import (
	"github.com/Lelisayohanes/erplite/backend/internal/shared/config"
	"github.com/labstack/echo/v5"
)

type App struct {
	Echo   *echo.Echo
	Config *config.Config
}

func New(cfg *config.Config) *App {
	e := echo.New()

	app := &App{
		Echo:   e,
		Config: cfg,
	}

	app.registerMiddleware()
	app.registerRoutes()

	return app
}

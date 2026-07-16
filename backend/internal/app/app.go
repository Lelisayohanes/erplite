package app

import (
	"github.com/Lelisayohanes/erplite/backend/internal/shared/config"
	"github.com/labstack/echo/v5"
)

type App struct {
	server *echo.Echo
	cfg    *config.Config
}

func New(cfg *config.Config) *App {
	e := echo.New()

	app := &App{
		server: e,
		cfg:    cfg,
	}

	app.registerMiddleware()
	app.registerRoutes()

	return app
}

func (a *App) Run() error {
	addr := a.cfg.Server.Host + ":" + a.cfg.Server.Port
	return a.server.Start(addr)
}

package app

import (
	"net/http"

	"github.com/Lelisayohanes/erplite/backend/internal/shared/config"
	"github.com/labstack/echo/v5"
)

type App struct {
	echoServer *echo.Echo
	cfg        *config.Config
	httpServer *http.Server
}

func New(cfg *config.Config) *App {

	e := echo.New()

	app := &App{
		echoServer: e,
		cfg:        cfg,
	}

	app.registerMiddleware()
	app.registerRoutes()
	app.createHTTPServer()

	return app
}

// Handler exposes the HTTP handler for testing.
func (a *App) Handler() http.Handler {
	return a.echoServer
}

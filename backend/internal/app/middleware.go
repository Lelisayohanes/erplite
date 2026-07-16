package app

import "github.com/labstack/echo/v5/middleware"

func (a *App) registerMiddleware() {
	a.server.Use(middleware.Recover())
	a.server.Use(middleware.RequestLogger())
}

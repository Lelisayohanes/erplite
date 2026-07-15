package app

import "github.com/labstack/echo/v5/middleware"

func (a *App) registerMiddleware() {
	a.Echo.Use(middleware.Recover())
	a.Echo.Use(middleware.RequestLogger())
}

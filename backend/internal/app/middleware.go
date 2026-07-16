package app

import "github.com/labstack/echo/v5/middleware"

func (a *App) registerMiddleware() {
	a.echoServer.Use(middleware.Recover())
	a.echoServer.Use(middleware.RequestLogger())
}

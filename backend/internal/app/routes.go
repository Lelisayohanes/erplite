package app

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (a *App) registerRoutes() {
	a.server.GET("/ping", ping)

	a.server.GET("/health", health)

	a.server.GET("/ready", ready)

	a.server.GET("/live", live)
}

func ping(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "pong",
	})
}

func health(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

func ready(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ready",
	})
}

func live(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "alive",
	})
}

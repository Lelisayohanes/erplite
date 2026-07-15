package app

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (a *App) registerRoutes() {
	a.Echo.GET("/ping", ping)

	a.Echo.GET("/health", health)

	a.Echo.GET("/ready", ready)

	a.Echo.GET("/live", live)
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

package main // Entry point of the Go application.

import (
	"net/http" // Provides HTTP status codes like http.StatusOK.

	"github.com/labstack/echo/v5"            // Echo web framework.
	"github.com/labstack/echo/v5/middleware" // Built-in middleware (logging, recovery, etc.).
)

func main() {
	// Create a new Echo application instance.
	e := echo.New()

	// Log every incoming HTTP request.
	e.Use(middleware.RequestLogger())

	// Recover from panics so the server doesn't crash.
	e.Use(middleware.Recover())

	// Register a GET endpoint for "/".
	e.GET("/ping", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
	})

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

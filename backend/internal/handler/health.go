package handler

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Health is the liveness handler — confirms the process is running.
func Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// Readiness is the readiness handler — confirms the process and all
// required dependencies (e.g. database) are available.
func Readiness(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		if db == nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "not_ready",
				"reason": "database unavailable",
			})
		}
		if err := db.PingContext(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "not_ready",
				"reason": "database unavailable",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "ready",
		})
	}
}

package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/labstack/echo/v4"
)

// CorrelationID extracts or generates a correlation identifier for request tracing.
// The ID is taken from the X-Correlation-ID header if present, otherwise a new
// UUID is generated. It is set on the response header and stored in the echo context.
func CorrelationID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get(echo.HeaderXCorrelationID)
			if id == "" {
				id = c.Response().Header().Get(echo.HeaderXRequestID)
			}
			if id == "" {
				id = generateID()
			}

			c.Set("correlation_id", id)
			c.Response().Header().Set(echo.HeaderXCorrelationID, id)

			return next(c)
		}
	}
}

// generateID returns a random 16-byte hex string for correlation tracing.
func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// RequestLogger logs every HTTP request as a structured entry including
// method, URI, status, latency, and the correlation identifier.
// Sensitive information (headers, body, auth tokens) is excluded.
func RequestLogger(log echo.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			latency := time.Since(start)
			status := c.Response().Status

			correlationID, _ := c.Get("correlation_id").(string)

			log.Infof(
				"correlation_id=%s method=%s uri=%s status=%d latency=%s remote_ip=%s",
				correlationID,
				c.Request().Method,
				c.Request().RequestURI,
				status,
				latency,
				c.RealIP(),
			)

			return err
		}
	}
}

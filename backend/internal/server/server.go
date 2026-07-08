package server

import (
	"erplite/backend/internal/container"
	"erplite/backend/internal/handler"
	"erplite/backend/internal/middleware"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
	ctr  *container.Container
}

func New(ctr *container.Container) *Server {
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.CorrelationID())
	e.Use(middleware.RequestLogger(e.Logger))

	// Routes
	e.GET("/healthz", handler.Health)
	e.GET("/readyz", handler.Readiness(ctr.DB))

	return &Server{echo: e, ctr: ctr}
}

func (s *Server) Start() error {
	return s.echo.Start(":" + s.ctr.Cfg.App.Port)
}

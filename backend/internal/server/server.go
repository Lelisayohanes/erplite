package server

import (
	"erplite/backend/internal/config"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
	cfg  *config.Config
}

func New(cfg *config.Config) *Server {
	e := echo.New()
	return &Server{echo: e, cfg: cfg}
}

func (s *Server) Start() error {
	return s.echo.Start(":" + s.cfg.App.Port)
}

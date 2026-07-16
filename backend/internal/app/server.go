package app

import (
	"net/http"
)

func (a *App) createHTTPServer() {

	a.httpServer = &http.Server{
		Addr: a.cfg.HTTP.Host + ":" + a.cfg.HTTP.Port,

		Handler: a.echoServer,
	}
}

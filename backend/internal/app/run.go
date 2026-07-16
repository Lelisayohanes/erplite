package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (a *App) Run() error {
	serverErr := make(chan error, 1)

	if err := a.startServer(serverErr); err != nil {
		return err
	}

	return a.waitForShutdown(serverErr)
}

func (a *App) startServer(serverErr chan<- error) error {
	addr := a.cfg.HTTP.Host + ":" + a.cfg.HTTP.Port

	a.httpServer = &http.Server{
		Addr:    addr,
		Handler: a.echoServer,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	return nil
}

func (a *App) waitForShutdown(serverErr <-chan error) error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverErr:
		return err

	case <-ctx.Done():
		return a.shutdown()
	}
}

func (a *App) shutdown() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	return a.httpServer.Shutdown(ctx)
}

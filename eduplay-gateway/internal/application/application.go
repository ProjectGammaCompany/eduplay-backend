package application

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type App struct {
	log    *slog.Logger
	server *http.Server
}

func New(
	log *slog.Logger,
	address string,
	timeout time.Duration,
	router *chi.Mux,
) *App {
	srv := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  timeout,
	}

	return &App{
		log:    log,
		server: srv,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "server.Run"

	if err := a.server.ListenAndServe(); err != nil {
		a.log.Error("failed to start server", slog.String("op", op), slog.String("err", err.Error()))
		return err
	}

	a.log.With(slog.String("op", op)).Info("gateway server started", slog.String("addr", a.server.Addr))

	return nil
}

func (a *App) Stop() {
	const op = "server.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gateway server", slog.String("addr", a.server.Addr))

	err := a.server.Shutdown(context.Background())
	if err != nil {
		a.log.Error("failed to stop server", slog.String("op", op), slog.String("err", err.Error()))
	}
}

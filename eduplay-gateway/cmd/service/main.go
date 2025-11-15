package main

import (
	"context"
	"eduplay-gateway/internal/application"
	"eduplay-gateway/internal/config"
	"eduplay-gateway/internal/http/handlers/users/changePassword"
	"eduplay-gateway/internal/http/handlers/users/deleteAccount"
	"eduplay-gateway/internal/http/handlers/users/getUserData"
	"eduplay-gateway/internal/http/handlers/users/refresh"
	"eduplay-gateway/internal/http/handlers/users/signIn"
	signOutUser "eduplay-gateway/internal/http/handlers/users/signOut"
	"eduplay-gateway/internal/http/handlers/users/signUp"
	"eduplay-gateway/internal/http/handlers/users/updateUserData"
	mwCors "eduplay-gateway/internal/http/middleware/cors"
	mwLogger "eduplay-gateway/internal/http/middleware/logger"
	uClient "eduplay-gateway/internal/pkg/clients/users"
	"eduplay-gateway/internal/pkg/usecases/users"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type App struct {
	GatewayServer *application.App
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	userClient, err := uClient.New(context.Background(), log,
		cfg.Clients.Users.Address,
		cfg.Clients.Users.Timeout,
		cfg.Clients.Users.Retries)
	if err != nil {
		log.Error("failed to create users client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	// router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(mwCors.CorsMiddleware)

	router.Route("/auth", func(r chi.Router) {
		r.Post("/signUp", signUp.New(log, users.New(log, userClient)))
		r.Post("/signIn", signIn.New(log, users.New(log, userClient)))
		r.Post("/refresh", refresh.New(log, users.New(log, userClient)))
		r.Get("/userData", getUserData.New(log, users.New(log, userClient)))
		r.Put("/userData", updateUserData.New(log, users.New(log, userClient)))
		r.Delete("/", deleteAccount.New(log, users.New(log, userClient)))
		r.Post("/changePassword", changePassword.New(log, users.New(log, userClient)))
		r.Post("/signOut", signOutUser.New(log, users.New(log, userClient)))
	})

	server := application.New(log, cfg.Server.Address, cfg.Server.Timeout, router)

	app := &App{
		GatewayServer: server,
	}

	go func() {
		app.GatewayServer.MustRun()
	}()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GatewayServer.Stop()
	log.Info("Gracefully stopped")
}

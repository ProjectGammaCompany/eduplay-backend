package main

import (
	"context"
	"eduplay-user/internal/application"
	"eduplay-user/internal/config"
	db "eduplay-user/internal/pkg/storage/postgres"
	users "eduplay-user/internal/pkg/usecase/user"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type App struct {
	GRPCServer *application.App
}

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	storage, _ := db.New(context.Background(), cfg.StoragePath)

	usersService := users.New(log, storage, cfg.SecretKey)

	grpcApp := application.New(log, usersService, cfg.GRPC.Port)

	app := &App{
		grpcApp,
	}

	go func() {
		app.GRPCServer.MustRun()
	}()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GRPCServer.Stop()
	log.Info("Gracefully stopped")
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

package main

import (
	"context"
	"eduplay-notification/internal/application"
	"eduplay-notification/internal/config"
	db "eduplay-notification/internal/pkg/storage"
	notifs "eduplay-notification/internal/pkg/usecase/notification"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
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

	storage, _ := db.New(context.Background(), cfg.StoragePath, cfg.EventStoragePath, cfg.UserStoragePath)

	usersService := notifs.New(log, storage, cfg.SecretKey)

	c := cron.New()
	_, err := c.AddFunc("@every 1m", func() {
		if err := storage.GetUserNotifications(context.Background()); err != nil {
			log.Error("failed to get user notifications", slog.String("error", err.Error()))
		}
	})
	if err != nil {
		log.Error("failed to add cron job", slog.String("error", err.Error()))
	}

	grpcApp := application.New(log, usersService, cfg.GRPC.Port)

	app := &App{
		grpcApp,
	}

	go func() {
		app.GRPCServer.MustRun()
	}()

	c.Start()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	c.Stop()

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

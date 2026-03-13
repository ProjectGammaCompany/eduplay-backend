package main

import (
	"context"
	"eduplay-event/internal/application"
	"eduplay-event/internal/config"
	db "eduplay-event/internal/pkg/storage/postgres"
	events "eduplay-event/internal/pkg/usecase/events"
	rabbit "eduplay-event/internal/pkg/usecase/rabbitmq"
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

	storage, err := db.New(context.Background(), cfg.StoragePath)
	if err != nil {
		log.Error("failed to create storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	rabbitMQ, err := rabbit.NewRabbitMQ(cfg.RabbitMQ, log, storage)
	if err != nil {
		log.Error("failed to create rabbitmq", slog.String("error", err.Error()))
		os.Exit(1)
	}

	c := cron.New()
	_, err = c.AddFunc("@hourly", func() {
		if err := storage.DeleteExpiredJoinCodes(context.Background()); err != nil {
			log.Error("failed to delete expired join codes", slog.String("error", err.Error()))
		}
	})
	if err != nil {
		log.Error("failed to add cron job", slog.String("error", err.Error()))
	}

	eventsService := events.New(log, storage, cfg.SecretKey)

	grpcApp := application.New(log, eventsService, cfg.GRPC.Port)

	app := &App{
		grpcApp,
	}

	go func() {
		app.GRPCServer.MustRun()
	}()

	go func() {
		rabbitMQ.ReceiveUserDeletedMessage()
	}()

	c.Start()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	c.Stop()

	app.GRPCServer.Stop()

	if err := rabbitMQ.Close(); err != nil {
		log.Error("failed to close rabbitmq", slog.String("error", err.Error()))
	}
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

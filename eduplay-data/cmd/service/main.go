package main

import (
	"context"
	"eduplay-data/internal/application"
	"eduplay-data/internal/config"
	evClient "eduplay-data/internal/pkg/clients/event"
	db "eduplay-data/internal/pkg/storage/postgres"
	"eduplay-data/internal/pkg/usecase/data"

	// rabbit "eduplay-event/internal/pkg/usecase/rabbitmq"
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

	storage, err := db.New(context.Background(), cfg.StoragePath)
	if err != nil {
		log.Error("failed to create storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// rabbitMQ, err := rabbit.NewRabbitMQ(cfg.RabbitMQ, log, storage)
	// if err != nil {
	// 	log.Error("failed to create rabbitmq", slog.String("error", err.Error()))
	// 	os.Exit(1)
	// }

	eventClient, err := evClient.New(context.Background(), log,
		cfg.Clients.Events.Address,
		cfg.Clients.Events.Timeout,
		cfg.Clients.Events.Retries)
	if err != nil {
		log.Error("failed to create events client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	dataService := data.New(log, storage, eventClient, cfg.SecretKey)

	grpcApp := application.New(log, dataService, cfg.GRPC.Port)

	app := &App{
		grpcApp,
	}

	go func() {
		app.GRPCServer.MustRun()
	}()

	// go func() {
	// 	rabbitMQ.ReceiveUserDeletedMessage()
	// }()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	// if err := rabbitMQ.Close(); err != nil {
	// 	log.Error("failed to close rabbitmq", slog.String("error", err.Error()))
	// }

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

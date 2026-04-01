package main

import (
	"eduplay-gateway/internal/application"
	"eduplay-gateway/internal/config"

	mwCors "eduplay-gateway/internal/http/middleware/cors"
	mwLogger "eduplay-gateway/internal/http/middleware/logger"
	"eduplay-gateway/internal/http/routers"
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

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	// router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(mwCors.CorsMiddleware)

	routers.UserRouter(router, log, cfg)
	routers.EventRouter(router, log, cfg)
	routers.DataRouter(router, log, cfg)
	routers.NotificationRouter(router, log, cfg)

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

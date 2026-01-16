package routers

import (
	"context"
	"eduplay-gateway/internal/config"
	"eduplay-gateway/internal/http/handlers/file/postFile"
	"log/slog"
	"os"

	evClient "eduplay-gateway/internal/pkg/clients/event"
	events "eduplay-gateway/internal/pkg/usecases/event"

	"github.com/go-chi/chi/v5"
)

func EventRouter(router chi.Router, log *slog.Logger, cfg *config.Config) chi.Router {
	eventClient, err := evClient.New(context.Background(), log,
		cfg.Clients.Events.Address,
		cfg.Clients.Events.Timeout,
		cfg.Clients.Events.Retries)
	if err != nil {
		log.Error("failed to create events client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router.Route("/", func(r chi.Router) {
		r.Post("/file", postFile.New(log, events.New(log, eventClient)))
	})

	return router
}

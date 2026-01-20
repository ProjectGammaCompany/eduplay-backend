package routers

import (
	"context"
	"eduplay-gateway/internal/config"
	"eduplay-gateway/internal/http/handlers/event/getEvent"
	"eduplay-gateway/internal/http/handlers/event/getEventBlocks"
	"eduplay-gateway/internal/http/handlers/event/getEventRole"
	"eduplay-gateway/internal/http/handlers/event/getEventSettings"
	"eduplay-gateway/internal/http/handlers/event/getPublicEvents"
	"eduplay-gateway/internal/http/handlers/event/getUserFavorites"
	"eduplay-gateway/internal/http/handlers/event/postEvent"
	"eduplay-gateway/internal/http/handlers/event/postEventBlock"
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
		r.Get("/events", getPublicEvents.New(log, events.New(log, eventClient)))
		r.Get("/events/personal/favorites", getUserFavorites.New(log, events.New(log, eventClient)))
	})

	router.Route("/event", func(r chi.Router) {
		r.Post("/", postEvent.New(log, events.New(log, eventClient)))
		r.Get("/{eventId}", getEvent.New(log, events.New(log, eventClient)))
		r.Get("/{eventId}/role", getEventRole.New(log, events.New(log, eventClient)))
		r.Get("/{eventId}/settings", getEventSettings.New(log, events.New(log, eventClient)))
		r.Post("/{eventId}/block", postEventBlock.New(log, events.New(log, eventClient)))
		r.Get("/{eventId}", getEventBlocks.New(log, events.New(log, eventClient))) // TODO check
	})
	return router
}

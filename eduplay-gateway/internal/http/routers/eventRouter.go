package routers

import (
	"context"
	"eduplay-gateway/internal/config"
	"eduplay-gateway/internal/http/handlers/event/deleteBlockById"
	"eduplay-gateway/internal/http/handlers/event/deleteBlockCondition"
	"eduplay-gateway/internal/http/handlers/event/deleteEventById"
	"eduplay-gateway/internal/http/handlers/event/deleteTaskById"
	"eduplay-gateway/internal/http/handlers/event/getAllTags"
	"eduplay-gateway/internal/http/handlers/event/getBlockConditions"
	"eduplay-gateway/internal/http/handlers/event/getBlockInfo"
	"eduplay-gateway/internal/http/handlers/event/getBlockTasks"
	"eduplay-gateway/internal/http/handlers/event/getBlocksForConditions"
	"eduplay-gateway/internal/http/handlers/event/getEvent"
	"eduplay-gateway/internal/http/handlers/event/getEventBlocks"
	"eduplay-gateway/internal/http/handlers/event/getEventPlayerInfo"
	"eduplay-gateway/internal/http/handlers/event/getEventRole"
	"eduplay-gateway/internal/http/handlers/event/getEventSettings"
	"eduplay-gateway/internal/http/handlers/event/getHistory"
	"eduplay-gateway/internal/http/handlers/event/getNextStage"
	"eduplay-gateway/internal/http/handlers/event/getOwnedEvents"
	"eduplay-gateway/internal/http/handlers/event/getPublicEvents"
	"eduplay-gateway/internal/http/handlers/event/getTaskById"
	"eduplay-gateway/internal/http/handlers/event/getUserFavorites"
	"eduplay-gateway/internal/http/handlers/event/postAnswer"
	"eduplay-gateway/internal/http/handlers/event/postBlockCondition"
	"eduplay-gateway/internal/http/handlers/event/postEvent"
	"eduplay-gateway/internal/http/handlers/event/postEventBlock"
	"eduplay-gateway/internal/http/handlers/event/postTask"
	"eduplay-gateway/internal/http/handlers/event/putBlockCondition"
	"eduplay-gateway/internal/http/handlers/event/putEvent"
	"eduplay-gateway/internal/http/handlers/event/putEventBlock"
	"eduplay-gateway/internal/http/handlers/event/putFavorite"
	"eduplay-gateway/internal/http/handlers/event/putNextStage"
	"eduplay-gateway/internal/http/handlers/event/putTask"
	"eduplay-gateway/internal/http/handlers/event/putTimestamp"
	"eduplay-gateway/internal/http/handlers/file/postFile"

	"log/slog"
	"os"

	evClient "eduplay-gateway/internal/pkg/clients/event"
	uClient "eduplay-gateway/internal/pkg/clients/user"
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

	userClient, err := uClient.New(context.Background(), log,
		cfg.Clients.Users.Address,
		cfg.Clients.Users.Timeout,
		cfg.Clients.Users.Retries)
	if err != nil {
		log.Error("failed to create users client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router.Route("/", func(r chi.Router) {
		r.Post("/file", postFile.New(log, events.New(log, eventClient, userClient)))
		r.Get("/events", getPublicEvents.New(log, events.New(log, eventClient, userClient)))
		r.Get("/events/personal/favorites", getUserFavorites.New(log, events.New(log, eventClient, userClient)))
		r.Put("/events/personal/favorites", putFavorite.New(log, events.New(log, eventClient, userClient)))
		r.Get("/events/personal/owned", getOwnedEvents.New(log, events.New(log, eventClient, userClient)))
		r.Get("/events/personal/created", getOwnedEvents.New(log, events.New(log, eventClient, userClient)))
		r.Get("/events/personal/history", getHistory.New(log, events.New(log, eventClient, userClient)))
		r.Get("/tags", getAllTags.New(log, events.New(log, eventClient, userClient)))
	})

	router.Route("/event", func(r chi.Router) {
		r.Post("/", postEvent.New(log, events.New(log, eventClient, userClient)))
		r.Get("/{eventId}", getEvent.New(log, events.New(log, eventClient, userClient)))
		r.Put("/{eventId}", putEvent.New(log, events.New(log, eventClient, userClient)))
		r.Get("/{eventId}/role", getEventRole.New(log, events.New(log, eventClient, userClient)))
		r.Get("/{eventId}/settings", getEventSettings.New(log, events.New(log, eventClient, userClient)))
		r.Post("/{eventId}/block", postEventBlock.New(log, events.New(log, eventClient, userClient)))
		r.Put("/{eventId}/blocks/{blockId}", putEventBlock.New(log, events.New(log, eventClient, userClient)))
		r.Get("/{eventId}", getEventBlocks.New(log, events.New(log, eventClient, userClient)))
		r.Post("/{eventId}/blocks/{blockId}/task", postTask.New(log, events.New(log, eventClient, userClient)))
		r.Put("/{eventId}/blocks/{blockId}/tasks/{taskId}", putTask.New(log, events.New(log, eventClient, userClient)))
		r.Post("/{eventId}/blocks/{blockId}/conditions", postBlockCondition.New(log, events.New(log, eventClient, userClient)))
		r.Put("/{eventId}/blocks/{blockId}/conditions/{conditionId}", putBlockCondition.New(log, events.New(log, eventClient, userClient)))
		r.Get("/{eventId}/blocks/{blockId}", getBlockInfo.New(log, events.New(log, eventClient, userClient)))
		r.Get("/{eventId}/blocks/{blockId}/conditions", getBlockConditions.New(log, events.New(log, eventClient, userClient)))
		r.Delete("/{eventId}/blocks/{blockId}/conditions/{conditionId}", deleteBlockCondition.New(log, events.New(log, eventClient, userClient)))
		r.Get("/{eventId}/blocks", getBlocksForConditions.New(log, events.New(log, eventClient, userClient)))
		r.Get("/{eventId}/blocks/{blockId}/tasks", getBlockTasks.New(log, events.New(log, eventClient, userClient)))
		r.Get("/{eventId}/blocks/{blockId}/tasks/{taskId}", getTaskById.New(log, events.New(log, eventClient, userClient)))
		r.Delete("/{eventId}/blocks/{blockId}/tasks/{taskId}", deleteTaskById.New(log, events.New(log, eventClient, userClient)))
		r.Delete("/{eventId}/blocks/{blockId}", deleteBlockById.New(log, events.New(log, eventClient, userClient)))
		r.Delete("/{eventId}", deleteEventById.New(log, events.New(log, eventClient, userClient)))
		r.Post("/{eventId}/blocks/{blockId}/tasks/{taskId}/answer", postAnswer.New(log, events.New(log, eventClient, userClient)))
		r.Get("/{eventId}/playerInfo", getEventPlayerInfo.New(log, events.New(log, eventClient, userClient)))
		r.Put("/{eventId}/blocks/{blockId}/tasks/{taskId}/timestamp", putTimestamp.New(log, events.New(log, eventClient, userClient)))
		r.Put("/{eventId}/nextStage", putNextStage.New(log, events.New(log, eventClient, userClient)))
		r.Get("/{eventId}/nextStage", getNextStage.New(log, events.New(log, eventClient, userClient)))
	})

	return router
}

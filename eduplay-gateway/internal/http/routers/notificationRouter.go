package routers

import (
	"context"
	"eduplay-gateway/internal/config"
	deleteNotifications "eduplay-gateway/internal/http/handlers/notification/DeleteNotification"
	getNotifications "eduplay-gateway/internal/http/handlers/notification/GetNotifications"
	"log/slog"
	"os"

	notifClient "eduplay-gateway/internal/pkg/clients/notification"
	notifs "eduplay-gateway/internal/pkg/usecases/notification"

	"github.com/go-chi/chi/v5"
)

func NotificationRouter(router chi.Router, log *slog.Logger, cfg *config.Config) chi.Router {
	notificationClient, err := notifClient.New(context.Background(), log,
		cfg.Clients.Notification.Address,
		cfg.Clients.Notification.Timeout,
		cfg.Clients.Notification.Retries)
	if err != nil {
		log.Error("failed to create notification client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router.Route("/notifications", func(r chi.Router) {
		r.Get("/", getNotifications.New(log, notifs.New(log, notificationClient)))
		r.Delete("/{id}", deleteNotifications.New(log, notifs.New(log, notificationClient)))
	})

	return router
}

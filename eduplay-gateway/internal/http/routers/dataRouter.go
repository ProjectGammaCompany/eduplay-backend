package routers

import (
	"context"

	"log/slog"
	"os"

	"eduplay-gateway/internal/config"
	"eduplay-gateway/internal/http/handlers/data/getPublicEvents"
	dClient "eduplay-gateway/internal/pkg/clients/data"
	data "eduplay-gateway/internal/pkg/usecases/data"

	"github.com/go-chi/chi/v5"
)

func DataRouter(router chi.Router, log *slog.Logger, cfg *config.Config) chi.Router {
	dataClient, err := dClient.New(context.Background(), log,
		cfg.Clients.Data.Address,
		cfg.Clients.Data.Timeout,
		cfg.Clients.Data.Retries)
	if err != nil {
		log.Error("failed to create data client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router.Route("/data", func(r chi.Router) {
		r.Get("/getPublicEvents", getPublicEvents.New(log, data.New(log, dataClient)))
	})

	return router
}

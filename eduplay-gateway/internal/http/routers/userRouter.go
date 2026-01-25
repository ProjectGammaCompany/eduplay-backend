package routers

import (
	"context"
	"eduplay-gateway/internal/config"
	"eduplay-gateway/internal/http/handlers/user/changePassword"
	"eduplay-gateway/internal/http/handlers/user/deleteAccount"
	"eduplay-gateway/internal/http/handlers/user/getProfile"
	"eduplay-gateway/internal/http/handlers/user/refresh"
	"eduplay-gateway/internal/http/handlers/user/signIn"
	signOutUser "eduplay-gateway/internal/http/handlers/user/signOut"
	"eduplay-gateway/internal/http/handlers/user/signUp"
	"log/slog"
	"os"

	uClient "eduplay-gateway/internal/pkg/clients/user"
	users "eduplay-gateway/internal/pkg/usecases/user"

	"github.com/go-chi/chi/v5"
)

func UserRouter(router chi.Router, log *slog.Logger, cfg *config.Config) chi.Router {
	userClient, err := uClient.New(context.Background(), log,
		cfg.Clients.Users.Address,
		cfg.Clients.Users.Timeout,
		cfg.Clients.Users.Retries)
	if err != nil {
		log.Error("failed to create users client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", signUp.New(log, users.New(log, userClient)))
		r.Post("/login", signIn.New(log, users.New(log, userClient)))
		r.Post("/refresh", refresh.New(log, users.New(log, userClient)))
		// r.Get("/userData", getUserData.New(log, users.New(log, userClient)))
		// r.Put("/userData", updateUserData.New(log, users.New(log, userClient)))
		r.Delete("/", deleteAccount.New(log, users.New(log, userClient)))
		r.Post("/changePassword", changePassword.New(log, users.New(log, userClient)))
		r.Put("/logout", signOutUser.New(log, users.New(log, userClient)))
	})

	router.Route("/profile", func(r chi.Router) {
		r.Get("/", getProfile.New(log, users.New(log, userClient)))
	})

	return router
}

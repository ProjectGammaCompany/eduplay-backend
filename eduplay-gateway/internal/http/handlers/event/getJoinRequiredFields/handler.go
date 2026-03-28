package getJoinRequiredFields

import (
	"context"
	"eduplay-gateway/internal/http/tokens"
	"eduplay-gateway/internal/lib"
	"eduplay-gateway/internal/storage"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UseCase interface {
	GetRole(ctx context.Context, userId string, eventId string) (int64, error)
	GetEventByJoinCode(ctx context.Context, joinCode string, userId string) (bool, error)
}

func New(log *slog.Logger, uc UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.event.getJoinRequiredFields"

		log = log.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())))

		accessToken := request.Header.Get("Authorization")
		if accessToken == "" {
			log.Error("no authorization token provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("no authorization token provided"))
			return
		}

		accessToken = strings.Split(request.Header.Get("Authorization"), " ")[1]
		if accessToken == "" {
			log.Error("no authorization token provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("user not authorized"))
			return
		}

		accessClaims, err := tokens.ValidateAccessToken(accessToken)
		if err != nil {
			if errors.Is(err, storage.ErrInvalidAccessToken) {
				log.Error("invalid access token", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusUnauthorized)
				render.JSON(writer, request, lib.Error("invalid access token"))
				return
			}
			if errors.Is(err, storage.ErrAccessTokenExpired) {
				log.Error("access token expired", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusUnauthorized)
				render.JSON(writer, request, lib.Error("access token expired"))
				return
			}
			log.Error("failed to validate access token", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("failed to validate access token"))
			return
		}

		joinCode := chi.URLParam(request, "joinCode")
		if joinCode == "" {
			log.Error("no join code provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("no join code provided"))
			return
		}

		if len(joinCode) != 6 {
			log.Error("invalid join code provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("invalid join code provided"))
			return
		}

		ret, err := uc.GetEventByJoinCode(request.Context(), joinCode, accessClaims.ID)

		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				log.Error(err.Error(), slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusNotFound)
				render.JSON(writer, request, err)
				return
			}
			if errors.Is(err, storage.ErrUserIsNotPlayer) {
				log.Error(err.Error(), slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusForbidden)
				render.JSON(writer, request, err)
				return
			}
			log.Error(err.Error(), slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, err)
			return
		}

		writer.WriteHeader(http.StatusOK)
		render.JSON(writer, request, map[string]bool{"groupFields": ret})
	}
}

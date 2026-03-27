package getJoinCode

import (
	"context"
	"eduplay-gateway/internal/http/tokens"
	"eduplay-gateway/internal/lib"
	"eduplay-gateway/internal/storage"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	eventModel "eduplay-gateway/internal/lib/models/event"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UseCase interface {
	GetRole(ctx context.Context, userId string, eventId string) (int64, error)
	GetJoinCode(ctx context.Context, eventId string) (*eventModel.JoinCode, error)
}

func New(log *slog.Logger, uc UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.event.getJoinCode"

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

		eventId := chi.URLParam(request, "eventId")
		if eventId == "" {
			log.Error("no Id provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("no Id provided"))
			return
		}

		isUUID := tokens.ValidateUUID(eventId)
		if !isUUID {
			log.Error("invalid id provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("invalid id provided"))
			return
		}

		role, err := uc.GetRole(request.Context(), accessClaims.ID, eventId)
		if err != nil {
			log.Error(err.Error(), slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, err)
			return
		}

		if role != 1 {
			log.Error("forbidden action")
			writer.WriteHeader(http.StatusForbidden)
			render.JSON(writer, request, lib.Error("user is forbidden to perform this action"))
			return
		}

		ret, err := uc.GetJoinCode(request.Context(), eventId)

		if err != nil {
			if errors.Is(err, storage.ErrJoinCodeRetryFailed) {
				log.Error(err.Error(), slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusNotFound)
				render.JSON(writer, request, err)
				return
			}
			log.Error(err.Error(), slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, err)
			return
		}

		writer.WriteHeader(http.StatusOK)
		render.JSON(writer, request, ret)
	}
}

package deleteNotifications

import (
	"context"
	"eduplay-gateway/internal/http/tokens"
	"eduplay-gateway/internal/lib"
	model "eduplay-gateway/internal/lib/models/notification"
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
	DeleteNotification(ctx context.Context, userId string, notificationId string) error
}

func New(log *slog.Logger, uc UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.notification.getNotifications"

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

		notificationId := chi.URLParam(request, "id")
		if notificationId == "" {
			log.Error("no event Id provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("no event Id provided"))
			return
		}

		isUUID := tokens.ValidateUUID(notificationId)
		if !isUUID {
			log.Error("invalid event id provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("invalid event id provided"))
			return
		}

		filters := &model.NotificationFilter{}

		filters.UserId = accessClaims.ID

		err = uc.DeleteNotification(request.Context(), accessClaims.ID, notificationId)
		if err != nil {
			log.Error(err.Error(), slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, err)
			return
		}

		writer.WriteHeader(http.StatusOK)
		render.JSON(writer, request, "notification "+notificationId+" deleted successfully")
	}
}

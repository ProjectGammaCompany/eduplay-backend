package postEventByJoinCode

import (
	"context"
	"eduplay-gateway/internal/http/tokens"
	"eduplay-gateway/internal/lib"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"eduplay-gateway/internal/storage"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UseCase interface {
	GetRole(ctx context.Context, userId string, eventId string) (int64, error)
	PostEventByJoinCode(ctx context.Context, participacionPwds *eventModel.ParticipationPasswords) (string, error)
}

func New(log *slog.Logger, uc UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.event.postEventByJoinCode"

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
			log.Error("no Id provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("no Id provided"))
			return
		}

		if len(joinCode) != 6 {
			log.Error("invalid join code provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("invalid join code provided"))
			return
		}

		var req *eventModel.ParticipationPasswords

		err = render.DecodeJSON(request.Body, &req)
		if err != nil {
			log.Error(storage.ErrInvalidRequest.Error(), slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, storage.ErrInvalidRequest)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			log.Error(storage.ErrValidationError.Error(), slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, storage.ErrValidationError.Error())
			return
		}

		req.UserId = accessClaims.ID
		req.JoinCode = joinCode

		ret, err := uc.PostEventByJoinCode(request.Context(), req)

		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				log.Error(err.Error(), slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusNotFound)
				render.JSON(writer, request, err.Error())
				return
			}
			if errors.Is(err, storage.ErrIncorrectPassword) {
				log.Error(err.Error(), slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusForbidden)
				render.JSON(writer, request, err.Error())
				return
			}
			if errors.Is(err, storage.ErrUserIsNotPlayer) {
				log.Error(err.Error(), slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusForbidden)
				render.JSON(writer, request, err.Error())
				return
			}
			if errors.Is(err, storage.ErrUserAlreadyExists) {
				log.Error(err.Error(), slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusNotFound)
				render.JSON(writer, request, err.Error())
				return
			}
			log.Error(err.Error(), slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, err)
			return
		}

		writer.WriteHeader(http.StatusOK)
		render.JSON(writer, request, map[string]string{"eventId": ret})
	}
}

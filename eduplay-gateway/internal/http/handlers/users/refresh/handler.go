package refresh

import (
	"context"
	"eduplay-gateway/internal/lib"
	model "eduplay-gateway/internal/lib/models/user"
	"eduplay-gateway/internal/storage"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UseCase interface {
	Refresh(ctx context.Context, tokens *model.RefreshToken) (*model.RefreshToken, error)
}

func New(log *slog.Logger, uc UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.users.refresh"

		log = log.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())))

		var req model.RefreshRequest

		// accessToken := request.Header.Get("Authorization")
		// if accessToken == "" {
		// 	log.Error("no authorization token provided")
		// 	writer.WriteHeader(http.StatusBadRequest)
		// 	render.JSON(writer, request, lib.Error("no authorization token provided"))
		// 	return
		// }

		err := render.DecodeJSON(request.Body, &req)
		if err != nil {
			log.Error("failed to deserialize request", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusUnauthorized)
			render.JSON(writer, request, lib.Error("fail to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			log.Error("fail to validate request", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusUnauthorized)
			render.JSON(writer, request, lib.ValidationError(validationErrors))
			return
		}

		tokens := &model.RefreshToken{
			// AccessToken:  accessToken,
			RefreshToken: req.RefreshToken,
		}

		newTokens, err := uc.Refresh(context.Background(), tokens)

		if err != nil {
			if errors.Is(err, storage.ErrInvalidRefreshToken) {
				log.Error("invalid refresh token", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusUnauthorized)
				render.JSON(writer, request, lib.Error("invalid refresh token"))
				return
			}
			if errors.Is(err, storage.ErrRefreshTokenExpired) {
				log.Error("refresh token expired", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusUnauthorized)
				render.JSON(writer, request, lib.Error("refresh token expired"))
				return
			}
			if errors.Is(err, storage.ErrRefreshTokenNotFound) {
				log.Error("refresh token not found", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusUnauthorized)
				render.JSON(writer, request, lib.Error("refresh token not found"))
				return
			}
			log.Error("failed to refresh token")
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("failed to refresh token"))
			return
		}

		log.Info("success to refresh token", slog.Any("response", newTokens))
		render.JSON(writer, request, newTokens)
	}
}

package getValidationCode

import (
	"context"
	"eduplay-gateway/internal/lib"
	model "eduplay-gateway/internal/lib/models/user"
	storage "eduplay-gateway/internal/storage"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UseCase interface {
	GetValidationCode(ctx context.Context, pd model.Code) (bool, error)
}

func New(log *slog.Logger, uc UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.users.GetValidationCode"

		log = log.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())))

		// accessToken := request.Header.Get("Authorization")
		// if accessToken == "" {
		// 	log.Error("no authorization token provided")
		// 	writer.WriteHeader(http.StatusBadRequest)
		// 	render.JSON(writer, request, lib.Error("no authorization token provided"))
		// 	return
		// }

		// accessToken = strings.Split(request.Header.Get("Authorization"), " ")[1]
		// if accessToken == "" {
		// 	log.Error("no authorization token provided")
		// 	writer.WriteHeader(http.StatusBadRequest)
		// 	render.JSON(writer, request, lib.Error("user not authorized"))
		// 	return
		// }

		// accessClaims, err := tokens.ValidateAccessToken(accessToken)
		// if err != nil {
		// 	if errors.Is(err, storage.ErrInvalidAccessToken) {
		// 		log.Error("invalid access token", slog.String("error", err.Error()))
		// 		writer.WriteHeader(http.StatusUnauthorized)
		// 		render.JSON(writer, request, lib.Error("invalid access token"))
		// 		return
		// 	}
		// 	if errors.Is(err, storage.ErrAccessTokenExpired) {
		// 		log.Error("access token expired", slog.String("error", err.Error()))
		// 		writer.WriteHeader(http.StatusUnauthorized)
		// 		render.JSON(writer, request, lib.Error("access token expired"))
		// 		return
		// 	}
		// 	log.Error("failed to validate access token", slog.String("error", err.Error()))
		// 	writer.WriteHeader(http.StatusInternalServerError)
		// 	render.JSON(writer, request, lib.Error("failed to validate access token"))
		// 	return
		// }

		var req model.Code

		// err := render.DecodeJSON(request.Body, &req)
		// if err != nil {
		// 	log.Error("failed to deserialize request", slog.String("error", err.Error()))
		// 	writer.WriteHeader(http.StatusBadRequest)
		// 	render.JSON(writer, request, lib.Error("fail to decode request"))
		// 	return
		// }

		// log.Info("request body decoded", slog.Any("request", req))

		// if err := validator.New().Struct(req); err != nil {
		// 	var validationErrors validator.ValidationErrors
		// 	errors.As(err, &validationErrors)
		// 	log.Error("fail to validate request", slog.String("error", err.Error()))
		// 	writer.WriteHeader(http.StatusBadRequest)
		// 	render.JSON(writer, request, lib.ValidationError(validationErrors))
		// 	return
		// }

		code := request.URL.Query().Get("code")
		if code == "" {
			log.Error("no code provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("no code provided"))
			return
		}

		req = model.Code{Code: code}

		validity, err := uc.GetValidationCode(context.Background(), req)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Error("user not found", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusForbidden)
				render.JSON(writer, request, lib.Error("user not found"))
				return
			}
			if errors.Is(err, storage.ErrCodeExpired) {
				log.Error("code expired", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusForbidden)
				render.JSON(writer, request, lib.Error("code expired"))
				return
			}
			log.Error("failed to send validation code")
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("failed to set new username"))
			return
		}

		log.Info("user validation code sent", slog.Any("response", req))
		render.JSON(writer, request, map[string]bool{"validity": validity})
	}
}

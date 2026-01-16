package signIn

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
	"github.com/go-playground/validator/v10"
)

type UseCase interface {
	SignIn(ctx context.Context, pd *model.UserPD) (*model.Credentials, error)
}

func New(log *slog.Logger, uc UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.users.signIn"

		log = log.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())))

		var req model.SignInRequest

		err := render.DecodeJSON(request.Body, &req)
		if err != nil {
			log.Error("failed to deserialize request", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("fail to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			log.Error("fail to validate request", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.ValidationError(validationErrors))
			return
		}

		pdModel := &model.UserPD{
			Email:    req.Email,
			Password: req.Password,
		}

		credentials, err := uc.SignIn(context.Background(), pdModel)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Error("user not found")
				writer.WriteHeader(http.StatusNotFound)
				render.JSON(writer, request, storage.ErrUserNotFound.Error())
				return
			}
			if errors.Is(err, storage.ErrIncorrectPassword) {
				log.Error("invalid password")
				writer.WriteHeader(http.StatusUnauthorized)
				render.JSON(writer, request, storage.ErrIncorrectPassword.Error())
				return
			}
			if errors.Is(err, storage.ErrIsActive) {
				log.Error("user is already active")
				writer.WriteHeader(http.StatusForbidden)
				render.JSON(writer, request, storage.ErrIsActive.Error())
				return
			}
			log.Error("failed to get credentials")
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("failed to get credentials"))
			return
		}

		log.Info("success to get credentials", slog.Any("response", credentials))
		render.JSON(writer, request, credentials)
	}
}

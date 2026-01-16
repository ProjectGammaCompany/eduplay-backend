package signUp

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
	SignUp(ctx context.Context, pd *model.SignUpRequest) (*model.Credentials, error)
}

func New(log *slog.Logger, uc UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.users.signUp"

		log = log.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())))

		var req model.SignUpRequest

		err := render.DecodeJSON(request.Body, &req)
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

		if req.Password != req.RepeatPassword {
			log.Debug("Passwords don't match")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, storage.ErrPasswordsNotMatch.Error())
			return
		}

		credentials, err := uc.SignUp(context.Background(), &req)
		if err != nil {
			if errors.Is(err, storage.ErrUserAlreadyExists) {
				writer.WriteHeader(http.StatusConflict)
				render.JSON(writer, request, storage.ErrUserAlreadyExists.Error())
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

package getPublicEvents

import (
	"context"
	"eduplay-gateway/internal/http/tokens"
	"eduplay-gateway/internal/lib"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"eduplay-gateway/internal/storage"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UseCase interface {
	GetPublicEvents(ctx context.Context, filters *eventModel.EventBaseFilters) (*eventModel.GetPublicEventsOut, error)
}

func New(log *slog.Logger, uc UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.event.getPublicEvents"

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

		// _, err := tokens.ValidateAccessToken(accessToken)
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

		filters := &eventModel.EventBaseFilters{}

		page := request.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}
		filters.Page, _ = strconv.ParseInt(page, 10, 64)

		maxOnPage := request.URL.Query().Get("maxOnPage")
		if maxOnPage == "" {
			maxOnPage = "10"
		}
		filters.MaxOnPage, _ = strconv.ParseInt(maxOnPage, 10, 64)
		// filters.Tags, _ = request.URL.Query()["tags"]
		// filters.DecliningRating, _ = strconv.ParseBool(request.URL.Query().Get("decliningRating"))
		// filters.Territorialized, _ = strconv.ParseBool(request.URL.Query().Get("territorialized"))
		// filters.Active, _ = strconv.ParseBool(request.URL.Query().Get("active"))
		// filters.UserId = accessClaims.ID

		filters.UserId = accessClaims.ID

		events, err := uc.GetPublicEvents(request.Context(), filters)
		if err != nil {
			log.Error("failed to get public events", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("failed to get public events"))
			return
		}

		writer.WriteHeader(http.StatusOK)
		render.JSON(writer, request, events)
	}
}

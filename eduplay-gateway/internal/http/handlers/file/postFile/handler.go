package postFile

import (
	"context"
	"eduplay-gateway/internal/http/tokens"
	"eduplay-gateway/internal/lib"
	"io"
	"log/slog"
	"os"
	"strings"

	"net/http"

	storage "eduplay-gateway/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type UseCase interface {
	SaveFile(ctx context.Context, fileName string, fileUUID string) (string, error)
}

func New(log *slog.Logger, uc UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.docs.postFile"

		log = log.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())))

		if request.Header.Get("Authorization") == "" {
			log.Error("no authorization token provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("no authorization token provided"))
			return
		}

		accessToken := strings.Split(request.Header.Get("Authorization"), " ")[1]
		if accessToken == "" {
			log.Error("no authorization token provided")
			writer.WriteHeader(http.StatusBadRequest)
			render.JSON(writer, request, lib.Error("no authorization token provided"))
			return
		}

		_, err := tokens.ValidateAccessToken(accessToken)
		if err != nil {
			if err == storage.ErrInvalidAccessToken {
				log.Error("invalid access token", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusForbidden)
				render.JSON(writer, request, lib.Error("invalid access token"))
				return
			}
			if err == storage.ErrAccessTokenExpired {
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

		err = request.ParseMultipartForm(64 << 20)
		if err != nil {
			log.Error("failed to parse multipart form", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("error parsing multipart form"))
			return
		}

		fileNames := []string{}

		// TODO in actuality works with only one file

		mForm := request.MultipartForm
		for k := range mForm.File {
			file, fileHeader, err := request.FormFile(k)
			if err != nil {
				log.Error("failed to get form file", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusInternalServerError)
				render.JSON(writer, request, lib.Error("error getting form file"))
				return
			}

			defer func() {
				if err := file.Close(); err != nil {
					log.Error("failed to close file", slog.String("error", err.Error()))
				}
			}()

			splitFileName := strings.Split(fileHeader.Filename, ".")

			newFileName := uuid.New().String() + "." + splitFileName[len(splitFileName)-1]

			// TODO Change path for server
			dst, err := os.Create("C:/Users/Cactus/goprojs/EduPlay-back/eduplay-backend/resources/" + newFileName)
			if err != nil {
				log.Error("failed to creating new file on server", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusInternalServerError)
				render.JSON(writer, request, lib.Error("error creating new file on server"))
				return
			}

			defer func() {
				if err := dst.Close(); err != nil {
					log.Error("failed to close destination file", slog.String("error", err.Error()))
				}
			}()

			if _, err := io.Copy(dst, file); err != nil {
				log.Error("failed to copy file to server", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusInternalServerError)
				render.JSON(writer, request, lib.Error("error copying file to server"))
				return
			}

			_, err = uc.SaveFile(context.Background(), fileHeader.Filename, newFileName)
			if err != nil {
				log.Error("failed to save file", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusInternalServerError)
				render.JSON(writer, request, lib.Error("error saving file"))
				return
			}

			fileNames = append(fileNames, newFileName)
		}

		log.Info("success to save file", slog.Any("response", fileNames[0]))
		render.JSON(writer, request, fileNames[0])
	}
}

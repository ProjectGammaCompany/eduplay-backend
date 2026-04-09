package getEventDownload

import (
	"bytes"
	"context"
	"eduplay-gateway/internal/http/tokens"
	"eduplay-gateway/internal/lib"
	eventModel "eduplay-gateway/internal/lib/models/event"
	"eduplay-gateway/internal/storage"

	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type UseCase interface {
	// GetRole(ctx context.Context, userId string, eventId string) (int64, error)
	SaveFile(ctx context.Context, fileName string, fileKey string, fileUUID string) (string, error)
	GetEventJson(ctx context.Context, eventId string) (*eventModel.EventDownloadFull, error)
	GetEvent(ctx context.Context, pd *eventModel.Id) (*eventModel.PostEventIn, error)
}

func New(log *slog.Logger, uc UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.event.getEventDownload"

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

		_, err := tokens.ValidateAccessToken(accessToken)
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

		allowDownloading, err := uc.GetEvent(request.Context(), &eventModel.Id{Id: eventId})
		if err != nil {
			log.Error("failed to get event", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("failed to get event"))
			return
		}
		if !allowDownloading.AllowDownloading {
			log.Error("event not allowed to download")
			writer.WriteHeader(http.StatusForbidden)
			render.JSON(writer, request, lib.Error("event not allowed to download"))
			return
		}

		eventJson, err := uc.GetEventJson(request.Context(), eventId)
		if err != nil {
			log.Error("failed to get event json", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("failed to get event json"))
			return
		}

		jsonData, err := json.MarshalIndent(eventJson, "", "  ")
		if err != nil {
			log.Error("failed to marshal json", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("failed to marshal json"))
			return
		}

		fileName := eventId + ".json"
		fileUUID := uuid.New().String()
		fileKey := fileUUID + ".json"

		var awsS3Client *s3.Client

		creds := credentials.NewStaticCredentialsProvider("minioadmin", "Determination1", "")

		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithCredentialsProvider(creds),
			config.WithRegion("us-east-1"))
		if err != nil {
			log.Error("failed to load aws config", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("error loading aws config"))
			return
		}

		endpoint := "http://minio.storage.svc.cluster.local:9000"

		awsS3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true
			o.BaseEndpoint = aws.String(endpoint)
			// o.EndpointResolver = s3.EndpointResolverFromURL(endpoint)
		})

		bucket := "eduplay-bucket"

		_, err = awsS3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{Bucket: aws.String(bucket)})
		if err != nil {
			_, err = awsS3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{Bucket: aws.String(bucket)})
			if err != nil {
				log.Error("failed to create bucket", slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusInternalServerError)
				render.JSON(writer, request, lib.Error("error creating bucket"))
				return
			}
		}

		out, err := awsS3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
		if err != nil {
			log.Error("minio list buckets failed", slog.String("err", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("error listing buckets"))
			return
		}
		log.Info("minio ok", slog.Int("buckets", len(out.Buckets)))

		//nolint:staticcheck // SA1019 this is intentional
		uploader := manager.NewUploader(awsS3Client)

		//nolint:staticcheck // SA1019 this is intentional
		result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fileKey),
			Body:   bytes.NewReader(jsonData),
		})

		if err != nil {
			log.Error("failed to upload file", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("error uploading file"))
			return
		}

		log.Info("file uploaded", slog.String("file_name", fileName), slog.String("file_path", result.Location))

		_, err = uc.SaveFile(context.Background(), fileName, fileKey, fileUUID)
		if err != nil {
			log.Error("failed to save file", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("error saving file"))
			return
		}

		log.Info("success to save file", slog.Any("response", fileKey))
		render.JSON(writer, request, map[string]string{"downloadPath": fileKey})
		// render.JSON(writer, request, eventJson)
	}
}

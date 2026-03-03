package postFile

import (
	"context"
	"eduplay-gateway/internal/http/tokens"
	"eduplay-gateway/internal/lib"
	"log/slog"
	"strings"

	"net/http"

	storage "eduplay-gateway/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type UseCase interface {
	SaveFile(ctx context.Context, fileName string, fileKey string, fileUUID string) (string, error)
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

		// TODO in actuality works with only one file

		log.Info("Headers:", slog.Any("headers", request.Header))
		log.Info("Body:", slog.Any("body", request.Body))

		// // Limiting request size to 10GB
		// request.Body = http.MaxBytesReader(writer, request.Body, 10<<30)

		// Parsing the multipart form data
		err = request.ParseMultipartForm(64 << 20)
		if err != nil {
			log.Error("failed to parse multipart form", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("error parsing multipart form"))
			return
		}

		file, handler, err := request.FormFile("file")
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

		fileName := handler.Filename

		splitFileName := strings.Split(fileName, ".")

		fileUUID := uuid.New().String()

		newFileName := fileUUID + "." + splitFileName[len(splitFileName)-1]

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
		key := "uploads/" + newFileName

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
			Key:    aws.String(key),
			Body:   file,
		})

		if err != nil {
			log.Error("failed to upload file", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("error uploading file"))
			return
		}

		log.Info("file uploaded", slog.String("file_name", fileName), slog.String("file_path", result.Location))

		_, err = uc.SaveFile(context.Background(), fileName, key, fileUUID)
		if err != nil {
			log.Error("failed to save file", slog.String("error", err.Error()))
			writer.WriteHeader(http.StatusInternalServerError)
			render.JSON(writer, request, lib.Error("error saving file"))
			return
		}

		log.Info("success to save file", slog.Any("response", newFileName))
		render.JSON(writer, request, newFileName)
	}
}

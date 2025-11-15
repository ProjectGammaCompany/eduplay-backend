package user

import (
	"context"
	dto "eduplay-user/internal/generated"
	"log/slog"
	"time"

	"eduplay-user/internal/model"
)

type storage interface {
	SignUpUser(ctx context.Context, name string, surname string, email string, organization string, phone string, password string) (string, error)
	SaveSession(ctx context.Context, userId string, refreshToken string, role string, accessLevel int64) error
	UpdateSession(ctx context.Context, userId string, refreshToken string) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetSessionByUserId(ctx context.Context, userId string) (*model.Session, error)
	GetUserByRefreshToken(ctx context.Context, refreshToken string) (*model.User, error)
	GetUserInfoByAuthToken(ctx context.Context, userID string) (*model.UserInfo, error)
	ChangeUserInfo(ctx context.Context, userInfo *dto.ChangeUserInfoIn, userID string) (*model.UserInfo, error)
	ChangeUserPassword(ctx context.Context, newHash, authToken string) error
	DeleteAccount(ctx context.Context, authToken string) error
	CheckRefreshTokenExists(ctx context.Context, refreshToken string) (bool, time.Time, error)
	GetUserPasswordById(ctx context.Context, userId string) (string, error)
	SignOutUser(ctx context.Context, userId string) error
}

type UseCase struct {
	log     *slog.Logger
	storage storage
	secret  string
}

func New(
	log *slog.Logger,
	st storage,
	secret string,
) *UseCase {
	return &UseCase{
		log:     log,
		storage: st,
		secret:  secret,
	}
}

package user

import (
	"context"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	dto "eduplay-user/internal/generated"
	"eduplay-user/internal/model"
)

func (a *UseCase) SignUpUser(ctx context.Context, in *dto.SignUpIn) (*model.Session, error) {
	const op = "sign_up_user.UseCase.SignUpUser"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("signing up user")

	passwordHash, err := HashPassword(in.Password)
	if err != nil {
		log.Error("failed to hash password", err.Error(), slog.String("password", in.Password))
		return &model.Session{}, err
	}

	userId, err := a.storage.SignUpUser(ctx, in.Email, string(passwordHash))
	if err != nil {
		log.Error("failed to sign up user", err.Error(), slog.String("email", in.Email))
		return &model.Session{}, err
	}

	// role := dto.Role_USER
	// modelRole := role.String()

	// userAccess := 0

	session, error := GenerateSession(ctx, userId, in.Email, a.secret)
	if error != nil {
		log.Error("failed to generate session", error.Error(), slog.String("userId", userId))
		return &model.Session{}, error
	}

	err = a.storage.SaveSession(ctx, userId, session.RefreshToken)
	if err != nil {
		log.Error("failed to save session", err.Error(), slog.String("userId", userId))
		return &model.Session{}, err
	}

	return session, nil
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

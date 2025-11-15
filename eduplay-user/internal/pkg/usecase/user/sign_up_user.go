package user

import (
	"context"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	dto "eduplay-user/internal/generated"
	"eduplay-user/internal/model"
)

func (a *UseCase) SignUpUser(ctx context.Context, name string, surname string, email string, organization string, phone string, password string) (*model.Session, error) {
	const op = "sign_up_user.UseCase.SignUpUser"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("signing up user")

	passwordHash, err := HashPassword(password)
	if err != nil {
		log.Error("failed to hash password", err.Error(), slog.String("password", password))
		return &model.Session{}, err
	}

	userId, err := a.storage.SignUpUser(ctx, name, surname, email, organization, phone, string(passwordHash))
	if err != nil {
		log.Error("failed to sign up user", err.Error(), slog.String("email", email))
		return &model.Session{}, err
	}

	role := dto.Role_USER
	modelRole := role.String()

	userAccess := 0

	session, error := GenerateSession(ctx, userId, name, surname, email, userAccess, modelRole, a.secret)
	if error != nil {
		log.Error("failed to generate session", error.Error(), slog.String("userId", userId))
		return &model.Session{}, error
	}

	err = a.storage.SaveSession(ctx, userId, session.RefreshToken, session.Role, int64(session.AccessLevel))
	if err != nil {
		log.Error("failed to save session", err.Error(), slog.String("userId", userId))
		return &model.Session{}, err
	}

	return session, nil
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

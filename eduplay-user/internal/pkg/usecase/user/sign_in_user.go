package user

import (
	"context"
	"eduplay-user/internal/model"
	"log/slog"

	errors "eduplay-user/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

func (a *UseCase) SignInUser(ctx context.Context, email string, password string) (*model.Session, error) {
	const op = "Users.SignInUser"

	log := a.log.With(
		slog.String("op", op),
	)

	currUser, err := a.storage.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if !CheckPasswordHash(password, currUser.Password) {
		return nil, errors.ErrIncorrectPassword
	}

	session, err := a.storage.GetSessionByUserId(ctx, currUser.Id)
	if session.IsActive {
		return nil, errors.ErrIsActive
	}
	if err != nil {
		return nil, err
	}

	newAccessToken, err := GenerateAuthToken(currUser.Id, currUser.Name, currUser.Surname, currUser.Email, session.AccessLevel, session.Role, a.secret)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	err = a.storage.UpdateSession(ctx, currUser.Id, newRefreshToken)
	if err != nil {
		return nil, err
	}

	session, err = a.storage.GetSessionByUserId(ctx, currUser.Id)
	if err != nil {
		return nil, err
	}

	session.AccessToken = newAccessToken

	log.Info("sign in user successfully")

	return session, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

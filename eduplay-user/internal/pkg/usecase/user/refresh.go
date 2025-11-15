package user

import (
	"context"
	"eduplay-user/internal/model"
	st "eduplay-user/internal/storage"
	"errors"
	"time"
)

func (a *UseCase) RefreshSession(ctx context.Context, refreshToken string) (*model.Tokens, error) {

	correctRefresh, expiry, err := a.storage.CheckRefreshTokenExists(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, st.ErrInvalidRefresh) {
			a.log.Error("refresh token not found in db")
			return &model.Tokens{Message: "refresh token not found in db"}, nil
		}
		return nil, err
	}
	if !correctRefresh {
		a.log.Error("invalid refresh token")
		return &model.Tokens{Message: "invalid refresh token"}, nil
	}
	if expiry.Before(time.Now().UTC().Add(3 * time.Hour)) {
		a.log.Error("refresh token expired {} < {}", expiry.String(), time.Now().UTC().Add(3*time.Hour).String())
		return &model.Tokens{Message: "refresh token expired"}, nil
	}

	currUser, err := a.storage.GetUserByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	currSession, err := a.storage.GetSessionByUserId(ctx, currUser.Id)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := GenerateAuthToken(currUser.Id, currUser.Name, currUser.Surname, currUser.Email, currSession.AccessLevel, currSession.Role, a.secret)
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

	return &model.Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

package user

import (
	"context"
	dto "eduplay-user/internal/generated"
	"eduplay-user/internal/model"
	"fmt"
	"log/slog"

	"github.com/dgrijalva/jwt-go"
)

func (a *UseCase) ChangeUserInfo(ctx context.Context, info *dto.ChangeUserInfoIn) (*model.UserInfo, error) {
	const op = "Users.ChangeUserInfo"

	log := a.log.With(
		slog.String("op", op),
	)

	var userId string

	token, err := jwt.Parse(info.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи")
		}
		return []byte(a.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok := claims["id"].(string); ok {
			userId = id
		} else {
			fmt.Println("Поле 'id' не найдено в токене")
		}
	} else {
		fmt.Println("Токен недействителен")
	}

	log.Info("attempting to get user info")

	currUser, err := a.storage.ChangeUserInfo(ctx, info, userId)
	if err != nil {
		return nil, err
	}

	return currUser, nil
}

package user

import (
	"context"
	"eduplay-user/internal/model"
	"fmt"
	"log/slog"

	"github.com/dgrijalva/jwt-go"
)

func (a *UseCase) GetUserInfo(ctx context.Context, accessToken string) (*model.UserInfo, error) {
	const op = "Users.GetUserInfo"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to get user info")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи")
		}
		return []byte(a.secret), nil
	})

	if err != nil {
		return nil, err
	}

	var userId string

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok := claims["id"].(string); ok {
			userId = id
		} else {
			fmt.Println("Поле 'id' не найдено в токене")
		}
	} else {
		fmt.Println("Токен недействителен")
	}

	currUser, err := a.storage.GetUserInfoByAuthToken(ctx, userId)
	if err != nil {
		return nil, err
	}

	return currUser, nil
}

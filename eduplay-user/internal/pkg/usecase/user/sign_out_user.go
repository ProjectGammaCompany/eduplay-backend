package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dgrijalva/jwt-go"
)

func (a *UseCase) SignOutUser(ctx context.Context, accessToken string) error {

	const op = "Users.SignOutUser"

	log := a.log.With(
		slog.String("op", op),
	)

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи")
		}
		return []byte(a.secret), nil
	})

	if err != nil {
		return err
	}

	var userId string

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok := claims["id"].(string); ok {
			userId = id
		} else {
			log.Warn("Поле 'id' не найдено в токене")
		}
	} else {
		log.Warn("Токен недействителен")
	}

	log.Info("attempting to sign out user")

	err = a.storage.SignOutUser(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

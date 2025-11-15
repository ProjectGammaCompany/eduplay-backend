package user

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log/slog"
)

func (a *UseCase) DeleteUserAccount(ctx context.Context, accessToken string) error {
	const op = "Users.DeleteUserAccount"

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
			fmt.Println("Поле 'id' не найдено в токене")
		}
	} else {
		fmt.Println("Токен недействителен")
	}

	log.Info("attempting to delete user account")

	err = a.storage.DeleteAccount(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

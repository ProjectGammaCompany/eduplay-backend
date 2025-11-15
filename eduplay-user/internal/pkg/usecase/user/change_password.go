package user

import (
	"context"
	dto "eduplay-user/internal/generated"
	"fmt"
	"log/slog"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (a *UseCase) ChangeUserPassword(ctx context.Context, in *dto.ChangePasswordIn) error {
	const op = "Users.ChangeUserPassword"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to change user password")

	var userId string

	token, err := jwt.Parse(in.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи")
		}
		return []byte(a.secret), nil
	})

	if err != nil {
		return err
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

	oldHash, err := a.storage.GetUserPasswordById(ctx, userId)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(in.Password))
	fmt.Println(err)
	if err != nil {
		return err
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", err.Error(), slog.String("password", in.NewPassword))
		return err
	}

	err = a.storage.ChangeUserPassword(ctx, string(newHash), userId)
	if err != nil {
		return err
	}

	return nil
}

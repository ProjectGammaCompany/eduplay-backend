package user

import (
	"context"
	dto "eduplay-user/internal/generated"
	"eduplay-user/internal/model"
	errs "eduplay-user/internal/storage"
	"errors"
	"log/slog"
)

func (a *UseCase) ChangeUserPassword(ctx context.Context, in *dto.ChangePasswordIn) (*model.Session, error) {
	const op = "Users.ChangeUserPassword"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to change user password")

	userId, email, err := a.storage.GetVerificationCode(ctx, in.Code)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return &model.Session{}, errs.ErrUserNotFound
		}
		return &model.Session{}, err
	}

	passwordHash, err := HashPassword(in.Password)
	if err != nil {
		log.Error("failed to hash password", err.Error(), slog.String("password", in.Password))
		return nil, err
	}

	err = a.storage.ChangeUserPassword(ctx, string(passwordHash), userId)
	if err != nil {
		return nil, err
	}

	session, error := GenerateSession(ctx, userId, email, a.secret)
	if error != nil {
		log.Error("failed to generate session", error.Error(), slog.String("userId", userId))
		return &model.Session{}, error
	}

	err = a.storage.SaveSession(ctx, userId, session.RefreshToken)
	if err != nil {
		log.Error("failed to save session", err.Error(), slog.String("userId", userId))
		return &model.Session{}, err
	}

	// var userId string

	// token, err := jwt.Parse(in.AccessToken, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("неверный метод подписи")
	// 	}
	// 	return []byte(a.secret), nil
	// })

	// if err != nil {
	// 	return err
	// }

	// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	// 	if id, ok := claims["id"].(string); ok {
	// 		userId = id
	// 	} else {
	// 		fmt.Println("Поле 'id' не найдено в токене")
	// 	}
	// } else {
	// 	fmt.Println("Токен недействителен")
	// }

	// oldHash, err := a.storage.GetUserPasswordById(ctx, userId)
	// if err != nil {
	// 	return err
	// }

	// err = bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(in.Password))
	// fmt.Println(err)
	// if err != nil {
	// 	return err
	// }

	// newHash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	// if err != nil {
	// 	log.Error("failed to hash password", err.Error(), slog.String("password", in.NewPassword))
	// 	return err
	// }

	return session, nil
}

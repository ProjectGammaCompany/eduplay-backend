package user

import (
	"context"
	"eduplay-user/internal/model"
	"eduplay-user/internal/rpc/converters"
	"fmt"

	dto "eduplay-user/internal/generated"

	"github.com/dgrijalva/jwt-go"
)

func (u *UseCase) GetUserAccess(ctx context.Context, authToken string) (*dto.GetUserAccessOut, error) {

	user, err := GetUserFromAccessToken(authToken, u.secret)
	if err != nil {
		return nil, err
	}

	session, err := u.storage.GetSessionByUserId(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	return &dto.GetUserAccessOut{Role: converters.StringToDto(session.Role), AccessLevel: int64(session.AccessLevel), UserId: user.Id}, nil
}

func GetUserFromAccessToken(accessToken string, secret string) (*model.User, error) {
	claims := jwt.MapClaims{
		"id":          "",
		"name":        "",
		"surname":     "",
		"email":       "",
		"accessLevel": 0,
		"role":        "",
		"exp":         0,
	}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	user := &model.User{
		Id:      claims["id"].(string),
		Name:    claims["name"].(string),
		Surname: claims["surname"].(string),
		Email:   claims["email"].(string),
	}

	return user, nil
}

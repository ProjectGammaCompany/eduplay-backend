package user

import (
	"context"
	"crypto/rand"
	"eduplay-user/internal/model"
	"encoding/base64"
	"io"

	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func GenerateAuthToken(userId string, email string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"id": userId,
		// "name":        name,
		// "surname":     surname,
		"email": email,
		// "accessLevel": accessLevel,
		// "role":        role,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//TODO: figure out where to store the key
	secretKey := []byte(secret)

	return token.SignedString(secretKey)
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func GenerateSession(ctx context.Context,
	userId string,
	email string,
	secret string) (*model.Session, error) {

	accessToken, err := GenerateAuthToken(userId, email, secret)
	if err != nil {
		//TODO: deal with error
		return &model.Session{}, err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		//TODO: deal with error
		return &model.Session{}, err
	}

	return &model.Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		// Role:         role,
		// AccessLevel:  accessLevel
	}, nil
}

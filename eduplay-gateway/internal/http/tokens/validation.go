package tokens

import (
	"eduplay-gateway/internal/storage"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// var secretKey = []byte("EduplaySecretkey")

type Claims struct {
	ID string `json:"id"`
	// Name        string `json:"name"`
	// Surname     string `json:"surname"`
	Email string `json:"email"`
	// AccessLevel int64  `json:"access_level"`
	// Role        string `json:"role"`
	Exp int64 `json:"exp"`
	jwt.StandardClaims
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	secretKey := []byte("EduplaySecretkey")

	fmt.Println(tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		if err.Error() == "Token is expired" {
			return nil, storage.ErrAccessTokenExpired
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	// data := claims["data"].(map[string]interface{})

	fmt.Println(claims)
	return &Claims{
		ID: claims["id"].(string),
		// Name:        claims["name"].(string),
		// Surname:     claims["surname"].(string),
		Email: claims["email"].(string),
		// AccessLevel: int64(claims["accessLevel"].(float64)),
		// Role:        claims["role"].(string),
		Exp: int64(claims["exp"].(float64))}, nil
}

func ValidateUUID(id string) bool {
	if err := uuid.Validate(id); err == nil {
		return true
	}

	return false
}

package core

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateAuthToken(userId uuid.UUID, jwtKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId.String(),
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

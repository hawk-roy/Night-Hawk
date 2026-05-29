package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("dev-secret-change-me")

func GenerateToken(userID int64, username string) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"iat":      now.Unix(),
		"exp":      now.Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

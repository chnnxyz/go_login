package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userId, username, roleName string, isSuper bool) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"role":     roleName,
		"isSuper":  isSuper,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

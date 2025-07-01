package helpers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(email string, pwReset bool) (string, error) {
	secret := []byte(os.Getenv("SERIALIZER_KEY"))

	// Use different salts (in Go, part of the claims)
	salt := os.Getenv("SERIALIZER_SALT")
	if pwReset {
		salt = os.Getenv("SERIALIZER_RESET_SALT")
	}

	// Embed salt in claims
	claims := jwt.MapClaims{
		"email": email,
		"salt":  salt,
		"exp":   time.Now().Add(30 * time.Minute).Unix(), // expire after 30 mins
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ConfirmToken(tokenStr string, pwReset bool) (string, error) {
	secret := []byte(os.Getenv("SERIALIZER_KEY"))
	expectedSalt := os.Getenv("SERIALIZER_SALT")
	if pwReset {
		expectedSalt = os.Getenv("SERIALIZER_RESET_SALT")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Optionally validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Validate expiration manually
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return "", errors.New("token expired")
			}
		}

		// Validate salt
		if salt, ok := claims["salt"].(string); !ok || salt != expectedSalt {
			return "", errors.New("invalid salt")
		}

		if email, ok := claims["email"].(string); ok {
			return email, nil
		}
	}

	return "", errors.New("invalid token")
}

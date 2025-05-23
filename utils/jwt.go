package utils

import (
    "os"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(username string, isSuper bool) (string, error) {
    secret := []byte(os.Getenv("JWT_SECRET"))

    claims := jwt.MapClaims{
        "username": username,
        "exp": time.Now().Add(time.Hour * 72).Unix(),
        "isSuper": isSuper,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secret)
}

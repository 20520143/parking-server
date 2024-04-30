package utils

import (
	jwt2 "github.com/golang-jwt/jwt/v4"
	"os"
)

func GenerateToken(userId string) (string, error) {
	claims := jwt2.MapClaims{}
	claims["id"] = userId
	claims["exp"] = EXPIRTE_TIME
	token := jwt2.NewWithClaims(jwt2.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv(JWT_SECRET_KEY)))
}

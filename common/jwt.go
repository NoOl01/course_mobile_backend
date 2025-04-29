package common

import (
	"backend_course/database"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateToken(account database.User) (string, string, error) {
	key := []byte(os.Getenv("JWT_SECRET"))

	accessTokenExpirationTime := time.Now().Add(1 * time.Hour).Unix()
	refreshTokenExpirationTime := time.Now().Add(30 * 24 * time.Hour).Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":  account.Id,
		"exp": accessTokenExpirationTime,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": refreshTokenExpirationTime,
	})

	accessTokenString, err := accessToken.SignedString(key)
	if err != nil {
		return "", "", err
	}
	refreshTokenString, err := refreshToken.SignedString(key)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

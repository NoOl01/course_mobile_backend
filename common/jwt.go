package common

import (
	"backend_course/database"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateToken(account *database.User) (string, string, error) {
	key := []byte(os.Getenv("JWT_SECRET"))

	accessTokenExpirationTime := time.Now().Add(1 * time.Hour).Unix()
	refreshTokenExpirationTime := time.Now().Add(30 * 24 * time.Hour).Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":  account.Id,
		"exp": accessTokenExpirationTime,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":  account.Id,
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

func Refresh(userId int64) (string, string, error) {
	key := []byte(os.Getenv("JWT_SECRET"))

	accessTokenExpirationTime := time.Now().Add(1 * time.Hour).Unix()
	refreshTokenExpirationTime := time.Now().Add(30 * 24 * time.Hour).Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":  userId,
		"exp": accessTokenExpirationTime,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":  userId,
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

func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	key := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func GetIdFromToken(claims jwt.MapClaims) string {
	id, idOk := claims["Id"].(string)
	if !idOk {
		return ""
	}
	return id
}

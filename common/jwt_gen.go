package common

import (
	"course_mobile/db_models"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"os"
)

func GenerateToken(account db_models.User, db *gorm.DB) (string, error) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	key := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":    account.Id,
		"Email": account.Email,
	})

	return token.SignedString(key)
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

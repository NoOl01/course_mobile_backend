package common

import (
	"backend_course/database"
	"backend_course/dto"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)

func Encrypt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	return string(hash)
}

func CheckPass(account dto.LoginDto, db *gorm.DB) (*database.User, error) {
	var user database.User

	err := db.Where("email = ?", account.Email).First(&user).Error
	if err != nil {
		fmt.Println(err)
		return nil, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(account.Password))
	if err != nil {
		fmt.Println(err)
		return nil, ErrWrongPassword
	}

	return &user, nil
}

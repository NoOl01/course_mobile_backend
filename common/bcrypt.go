package common

import (
	"course_mobile/db_models"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Encrypt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	return string(hash)
}

func CheckPass(account NewUser, db *gorm.DB) string {
	var user db_models.User

	err := db.Where("email = ?", account.Email).First(&user).Error
	if err != nil {
		fmt.Println(err)
		return "user not found"
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(account.Password))
	if err != nil {
		fmt.Println(err)
		return "wrong password"
	}

	return "Ok"
}

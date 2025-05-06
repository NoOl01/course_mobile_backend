package common

import (
	"backend_course/database"
	"errors"
	"gorm.io/gorm"
)

var (
	userNotFound        = errors.New("user not found")
	productNotFound     = "product not found"
	insufficientBalance = "user balance is insufficient"
)

func TryTransaction(db *gorm.DB, id, productId int64, product database.Product, order database.Order) error {
	var user database.User

	if err := checkTransaction(db, id, productId, &user, &product); err != nil {
		if errors.Is(err, productNotFound) {
			return err
		}

		if err := db.Model(&user).Where("id = ?", id).Update("balance", float64(user.Balance)-product.Price).Error; err != nil {
			return err
		}
		order.Status = ""
		if err := db.Create(&order).Error; err != nil {
			return err
		}
	}

	return nil
}

func checkTransaction(db *gorm.DB, id, productId int64, user *database.User, product *database.Product) error {
	if err := db.Model(&user).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userNotFound
		}
		return err
	}
	if err := db.Model(&product).Where("id = ?", productId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(productNotFound)
		}
		return err
	}

	if float64(user.Balance) < product.Price {
		return errors.New(insufficientBalance)
	}

	return nil
}

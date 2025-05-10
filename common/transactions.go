package common

import (
	"backend_course/database"
	"backend_course/statuses"
	"errors"
	"gorm.io/gorm"
)

func TryTransaction(db *gorm.DB, id, productId int64, product *database.Product, user *database.User) (string, error) {
	if err := checkTransaction(db, id, productId, user, product); err != nil {
		if errors.Is(err, statuses.InsufficientBalance) {
			return statuses.ResultAwaitingPayment, nil
		}

		return "", err
	}

	return statuses.ResultOk, nil
}

func checkTransaction(db *gorm.DB, id, productId int64, user *database.User, product *database.Product) error {
	if err := db.Where("id = ?", id).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return statuses.UserNotFound
		}
		return err
	}
	if err := db.Where("id = ?", productId).First(product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return statuses.ProductNotFound
		}
		return err
	}

	if user.Balance < product.Price {
		return statuses.InsufficientBalance
	}

	return nil
}

package common

import (
	"backend_course/database"
	"backend_course/statuses"
	"errors"
	"gorm.io/gorm"
)

func TryTransaction(db *gorm.DB, id, productId int64, product *database.Product, user *database.User, count int) (string, error) {
	if err := checkTransaction(db, id, productId, user, product, count); err != nil {
		if errors.Is(err, statuses.InsufficientBalance) {
			return statuses.ResultAwaitingPayment, nil
		}

		return "", err
	}

	return statuses.ResultPaid, nil
}

func checkTransaction(db *gorm.DB, id, productId int64, user *database.User, product *database.Product, count int) error {
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

	if user.Balance < (product.Price * float64(count)) {
		return statuses.InsufficientBalance
	}

	return nil
}

func RefundMoney(db *gorm.DB, id int64, order database.Order) error {
	switch order.Status {
	case statuses.ResultPaid:
		var user database.User
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			return err
		}
		user.Balance += order.Price
		if err := db.Save(&user).Error; err != nil {
			return err
		}
	}

	order.Status = statuses.ResultCanceled
	if err := db.Save(&order).Error; err != nil {
		return err
	}
	return nil
}

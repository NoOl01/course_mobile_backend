package statuses

import "errors"

// Errors
var (
	UserNotFound        = errors.New("user not found")
	ProductNotFound     = errors.New("product not found")
	InsufficientBalance = errors.New("user balance is insufficient")
)

// Results
const (
	ResultAwaitingPayment = "Ожидает оплаты"
	ResultOk              = "Оплачен"
)

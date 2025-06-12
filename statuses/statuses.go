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
	ResultPaid            = "Оплачен"
	ResultCanceled        = "Отменен"
	ResultOnTheWay        = "В пути"
	ResultSortingCenter   = "В сортировочном центре"
	ResultSentForDelivery = "Передано в доставку"
	ResultReceived        = "Получено"
)

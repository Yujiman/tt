package onlinepayment

type CheckPaymentRequest struct {
	OrderID uint `json:"order_id"`
}

type CheckPaymentResult struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
}

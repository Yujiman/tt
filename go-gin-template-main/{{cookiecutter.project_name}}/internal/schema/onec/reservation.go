package onec

type ReservationRequest struct {
	OrderID uint `json:"order_id"`
}

type ReservationResponse struct {
	ReceiptNumber string `json:"receipt_number"`
}

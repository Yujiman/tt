package onec

type GetOrderStatusRequest struct {
	OrderID uint `json:"order_id"`
}

type OrderStatusResponse struct {
	OrderID uint   `json:"order_id"`
	Status  string `json:"status"`
}

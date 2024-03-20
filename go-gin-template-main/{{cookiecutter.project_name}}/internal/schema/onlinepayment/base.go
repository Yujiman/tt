package onlinepayment

type Response[T any] struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Result  T      `json:"result"`
}

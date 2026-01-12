package accrual

type OrderResponse struct {
	Number  string  `json:"order"`
	Status  Status  `json:"status"`
	Accrual float64 `json:"accrual"`
}

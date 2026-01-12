package models

import "time"

type OrderStatus string

const OrderTable = "orders"

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

var OrderStatuses = []OrderStatus{OrderStatusNew, OrderStatusProcessing, OrderStatusInvalid, OrderStatusProcessed}

type Order struct {
	UploadedAt  *time.Time  `db:"uploaded_at"`
	ProcessedAt *time.Time  `db:"processed_at"`
	OrderNumber string      `db:"order_number"`
	Status      OrderStatus `db:"status"`
	UserID      uint        `db:"user_id"`
	Accrual     float64     `db:"accrual"`
}

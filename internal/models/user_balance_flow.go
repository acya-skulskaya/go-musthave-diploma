package models

import "time"

const UserBalanceFlowTable = "user_balance_flow"

type UserBalanceFlow struct {
	ProcessedAt *time.Time `db:"processed_at"`
	OrderNumber string     `db:"order_number"`
	UserID      uint       `db:"user_id"`
	Amount      float64    `db:"amount"`
}

package models

import "time"

const UserBalanceTable = "user_balance"

type UserBalance struct {
	UpdatedAt *time.Time `db:"updated_at"`
	UserID    uint       `db:"user_id"`
	Current   float64    `db:"current"`
	Withdrawn float64    `db:"withdrawn"`
}

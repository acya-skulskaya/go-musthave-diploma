package models

import "time"

const UserTable = "users"

type User struct {
	CreatedAt *time.Time `db:"created_at"`
	Login     string     `db:"login"`
	Password  string     `db:"password"`
	ID        uint       `db:"id"`
}

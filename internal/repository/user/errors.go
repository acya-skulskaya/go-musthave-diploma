package user

import "errors"

const (
	ErrMsgCouldNotCreate = "could not create user"
	ErrMsgCouldNotLogIn  = "could not log in"
	ErrMsgLoginNotFound  = "user with this login not found"
)

var ErrLoginAlreadyExists = errors.New("user with this login already exists")
var ErrNotFound = errors.New(ErrMsgLoginNotFound)

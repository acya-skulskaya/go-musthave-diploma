package order

import "errors"

const (
	ErrMsgCouldNotCreate = "could not create order"
)

var ErrOrderAlreadyExistsFromAnotherUser = errors.New("order number was already uploaded by another user")
var ErrOrderAlreadyExistsByCurrentUser = errors.New("order number by current user already exists")
var ErrOrderNumberAlreadyExists = errors.New("order number already exists")
var ErrNotFound = errors.New("order not found")

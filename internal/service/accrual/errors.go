package accrual

import "errors"

var ErrNoContent = errors.New("no content returned")
var ErrTooManyRequests = errors.New("too many requests")
var ErrUnknownOrderStatus = errors.New("unknown order status")

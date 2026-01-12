package balance

import "errors"

const (
	ErrMsgCouldNotCreate = "could not create order"
)

var ErrNotEnoughBalanceToWithdraw = errors.New("user does not have enough balance to withdraw")
var ErrWithdrawalForThisOrderAlreadyExists = errors.New("withdrawal for order already exists")
var ErrAccrualForThisOrderAlreadyExists = errors.New("accrual for order already exists")
var ErrNoWithdrawals = errors.New("user has no withdrawals")

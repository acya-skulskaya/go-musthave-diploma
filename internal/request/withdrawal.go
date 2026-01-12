package request

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/validator"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/validator/luhn"
	"strconv"
)

type Withdrawal struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (o Withdrawal) Valid() validator.Problems {
	problems := make(validator.Problems)

	id, err := strconv.Atoi(o.Order)
	if err != nil {
		problems[JSONFieldOrder] = ErrorMsgCouldNotConvertToInt
	} else if !luhn.Valid(id) {
		problems[JSONFieldOrder] = ErrorMsgInvalidLuhn
	}

	if o.Sum <= 0 {
		problems[JSONFieldSum] = ErrorMsgShouldBeMoreThanZero
	}

	return problems
}

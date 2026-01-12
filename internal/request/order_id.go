package request

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/validator"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/validator/luhn"
	"strconv"
)

type OrderID string

func (o OrderID) Valid() validator.Problems {
	problems := make(validator.Problems)

	id, err := strconv.Atoi(string(o))
	if err != nil {
		problems[JSONFieldBody] = ErrorMsgCouldNotConvertToInt
		return problems
	}

	if !luhn.Valid(id) {
		problems[JSONFieldBody] = ErrorMsgInvalidLuhn
	}

	return problems
}

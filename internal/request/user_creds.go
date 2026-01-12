package request

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/validator"
)

type UserCredits struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (o UserCredits) Valid() validator.Problems {
	problems := make(validator.Problems)

	if o.Login == "" {
		problems[JSONFieldLogin] = ErrorMsgEmptyLogin
	}

	if o.Password == "" {
		problems[JSONFieldPassword] = ErrorMsgEmptyPassword
	}

	return problems
}

package request

const (
	JSONFieldBody     = "body"
	JSONFieldLogin    = "login"
	JSONFieldPassword = "password"
	JSONFieldOrder    = "order"
	JSONFieldSum      = "sum"

	ErrorMsgCouldNotConvertToInt     = "could not convert value to int"
	ErrorMsgCouldNotConvertToFloat64 = "could not convert value to float64"
	ErrorMsgShouldBeMoreThanZero     = "should be more than zero"
	ErrorMsgInvalidLuhn              = "could not validate value with luhn"
	ErrorMsgEmptyLogin               = "login is required"
	ErrorMsgEmptyPassword            = "password is required"
)

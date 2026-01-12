package auth

import "errors"

var ErrTokenIsNotValid = errors.New("token is not valid")
var ErrCouldNotGetUserFromContext = errors.New("could not get usr from context")

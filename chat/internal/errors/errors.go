package chaterrors

import "errors"

var (
	UsernameAlreadyExists = errors.New("username already exists")
)

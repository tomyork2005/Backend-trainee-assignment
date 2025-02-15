package auth

import "errors"

var (
	errStorage         = errors.New("storage error")
	errInvalidPassword = errors.New("invalid password")
	errHashing         = errors.New("hashing error")
)

package autherrors

import "errors"

var (
	ErrStorage         = errors.New("storage error")
	ErrInvalidPassword = errors.New("invalid password")
	ErrHashing         = errors.New("hashing error")
)

package serverrors

import "errors"

var (
	ErrStorage     = errors.New("storage error")
	ErrInvalidData = errors.New("invalid data")
)

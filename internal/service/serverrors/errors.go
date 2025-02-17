package serverrors

import "errors"

var (
	ErrStorage          = errors.New("storage error")
	ErrInvalidTarget    = errors.New("target not found or invalid")
	ErrBalanceNotEnough = errors.New("balance not enough")
	ErrMerchDontExist   = errors.New("merch don't exist")
)

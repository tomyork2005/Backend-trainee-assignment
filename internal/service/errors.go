package service

import "errors"

var (
	ErrStorage                 = errors.New("storage error")
	ErrInvalidTarget           = errors.New("target not found or invalid")
	ErrBalanceNotEnough        = errors.New("balance not enough")
	ErrMerchDontExist          = errors.New("merch don't exist")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrTokenExpired            = errors.New("token expired")
	ErrInvalidToken            = errors.New("invalid token")
	ErrParsingToken            = errors.New("error with parse token")
	ErrHashing                 = errors.New("hashing error")
	ErrUnexpectedHashAlgorithm = errors.New("unexpected hash algorithm")
	ErrGenerateToken           = errors.New("error generating token")
)

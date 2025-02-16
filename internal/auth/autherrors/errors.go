package autherrors

import "errors"

var (
	ErrStorage                 = errors.New("storage error")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrInvalidToken            = errors.New("invalid token")
	ErrTokenExpired            = errors.New("token expired")
	ErrParsingToken            = errors.New("error with parse token")
	ErrHashing                 = errors.New("hashing error")
	ErrUnexpectedHashAlgorithm = errors.New("unexpected hash algorithm")
)

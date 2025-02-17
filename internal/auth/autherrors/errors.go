package autherrors

import "errors"

var (
	ErrStorage                 = errors.New("storage error")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrTokenExpired            = errors.New("token expired")
	ErrInvalidToken            = errors.New("invalid token")
	ErrParsingToken            = errors.New("error with parse token")
	ErrHashing                 = errors.New("hashing error")
	ErrUnexpectedHashAlgorithm = errors.New("unexpected hash algorithm")
	ErrGenerateToken           = errors.New("error generating token")
)

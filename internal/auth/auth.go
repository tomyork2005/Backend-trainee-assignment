package auth

import "time"

type authService struct {
	TokenTTL   time.Duration
	SigningKey string
}

func NewAuthService(tokenTTL time.Duration, signingKey string) *authService {
	return &authService{
		TokenTTL:   tokenTTL,
		SigningKey: signingKey,
	}
}

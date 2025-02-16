package auth

import "context"

func (auth *authService) ParseToken(ctx context.Context, token string) (string, error) {
	return token, nil
}

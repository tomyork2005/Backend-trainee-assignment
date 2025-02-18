package auth

import (
	"Backend-trainee-assignment/internal/service"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (auth *AuthService) ParseToken(ctx context.Context, token string) (string, error) {
	tkn, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w:%v", service.ErrUnexpectedHashAlgorithm, token.Header["alg"])
		}
		return []byte(auth.signingKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("%w:%w", service.ErrParsingToken, err)
	}

	if !tkn.Valid {
		return "", service.ErrInvalidToken
	}

	claims, ok := tkn.Claims.(*UserClaims)
	if !ok {
		return "", service.ErrInvalidToken
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return "", service.ErrTokenExpired
	}

	return claims.Username, nil
}

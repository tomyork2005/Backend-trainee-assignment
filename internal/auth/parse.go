package auth

import (
	"Backend-trainee-assignment/internal/auth/autherrors"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (auth *authService) ParseToken(ctx context.Context, token string) (string, error) {
	tkn, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w:%v", autherrors.ErrUnexpectedHashAlgorithm, token.Header["alg"])
		}
		return []byte(auth.signingKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("%w:%w", autherrors.ErrParsingToken, err)
	}

	if !tkn.Valid {
		return "", autherrors.ErrInvalidToken
	}

	claims, ok := tkn.Claims.(*UserClaims)
	if !ok {
		return "", autherrors.ErrInvalidToken
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return "", autherrors.ErrTokenExpired
	}

	return claims.Username, nil
}

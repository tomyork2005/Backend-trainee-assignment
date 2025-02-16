package auth

import (
	"Backend-trainee-assignment/internal/auth/autherrors"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

const (
	issuerName = "avitoMerchShop"
)

func generateToken(userID int64, tokenTTL time.Duration, signingKey string) (string, error) {

	claims := UserClaims{
		UserId: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuerName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		slog.Error("Error with signing token", slog.String("error", err.Error()))
		return "", err
	}

	return tokenString, nil
}

func (auth *authService) GetOrCreateTokenByCredentials(ctx context.Context, username, providedPassword string) (string, error) {

	user, err := auth.storage.GetUserIDByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("%w:%w", autherrors.ErrStorage, err)
	}

	// if user not found create new and return his token

	if user == nil {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(providedPassword), bcrypt.DefaultCost)
		if err != nil {
			return "", fmt.Errorf("%w:%w", autherrors.ErrHashing, err)
		}

		user, err = auth.storage.CreateUser(ctx, username, string(hashedPassword))
		if err != nil {
			return "", fmt.Errorf("%w:%w", autherrors.ErrStorage, err)
		}

		return generateToken(user.ID, auth.tokenTTL, auth.signingKey)
	}

	// compare passwords and return token

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return "", fmt.Errorf("%w:%w", autherrors.ErrInvalidPassword, err)
	}

	return generateToken(user.ID, auth.tokenTTL, auth.signingKey)
}

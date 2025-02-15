package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

func generateToken(userID int64, tokenTTL time.Duration, signingKey string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		jwtUserID: userID,
		jwtExp:    jwt.NewNumericDate(time.Now().Add(tokenTTL)),
	})

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		slog.Error("Error with signing token", slog.String("error", err.Error()))
		return "", err
	}

	return tokenString, nil
}

func (auth *authService) GetOrCreateTokenByCredentials(ctx context.Context, username, providedPassword string) (string, error) {

	user, err := auth.storage.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("%w:%w", errStorage, err)
	}

	// if user not found create new and return his token

	if user == nil {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(providedPassword), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("Error with hashing password", slog.String("error", err.Error()))
			return "", fmt.Errorf("%w:%w", errHashing, err)
		}

		user, err = auth.storage.CreateUser(ctx, username, string(hashedPassword))
		if err != nil {
			return "", fmt.Errorf("%w:%w", errStorage, err)
		}

		return generateToken(user.ID, auth.tokenTTL, auth.signingKey)
	}

	// compare passwords and return token

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return "", fmt.Errorf("%w:%w", errInvalidPassword, err)
	}

	return generateToken(user.ID, auth.tokenTTL, auth.signingKey)
}

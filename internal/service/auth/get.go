package auth

import (
	"Backend-trainee-assignment/internal/service"
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

func generateToken(username string, tokenTTL time.Duration, signingKey string) (string, error) {

	claims := UserClaims{
		Username: username,
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

func (auth *AuthService) GetOrCreateTokenByCredentials(ctx context.Context, username, providedPassword string) (string, error) {

	user, err := auth.storage.GetUserByUsername(ctx, username)
	if err != nil {
		slog.Error("Error getting user by username", slog.String("username", username), slog.String("error", err.Error()))

		return "", fmt.Errorf("%w:%w", service.ErrStorage, err)
	}

	// if user not found create new and return his token

	if user == nil {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(providedPassword), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("Error hashing password", slog.String("username", username))

			return "", fmt.Errorf("%w:%w", service.ErrHashing, err)
		}

		user, err = auth.storage.CreateUser(ctx, username, string(hashedPassword))
		if err != nil {
			slog.Error("Error creating user", slog.String("username", username), slog.String("error", err.Error()))
			return "", fmt.Errorf("%w:%w", service.ErrStorage, err)
		}

		token, err := generateToken(user.Username, auth.tokenTTL, auth.signingKey)
		if err != nil {
			slog.Error("Error generating token", slog.String("username", username))
			return "", fmt.Errorf("%w:%w", service.ErrGenerateToken, err)
		}

		slog.Info("Successfully generated token and create user", slog.String("username", username))

		return token, nil
	}

	// compare passwords and return token

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return "", fmt.Errorf("%w:%w", service.ErrInvalidPassword, err)
	}

	slog.Info("Successfully generated token", slog.String("username", username))

	return generateToken(user.Username, auth.tokenTTL, auth.signingKey)
}

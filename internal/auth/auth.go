package auth

import (
	"Backend-trainee-assignment/internal/model/service"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type storage interface {
	GetUserByUsername(ctx context.Context, username string) (*service.User, error)
	CreateUser(ctx context.Context, username, password string) (*service.User, error)
}
type authService struct {
	tokenTTL   time.Duration
	signingKey string

	storage storage
}

func NewAuthService(tokenTTL time.Duration, signingKey string, db storage) *authService {
	return &authService{
		tokenTTL:   tokenTTL,
		signingKey: signingKey,

		storage: db,
	}
}

type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

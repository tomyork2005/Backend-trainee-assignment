package storage

import (
	"Backend-trainee-assignment/internal/model"
	"context"
)

type storage struct {
}

func NewStorage() *storage {
	return &storage{}
}

func (s *storage) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return nil, nil
}

func (s *storage) CreateUser(ctx context.Context, username, password string) (*model.User, error) {
	return &model.User{
		ID:       123,
		Username: username,
		Password: password,
	}, nil
}

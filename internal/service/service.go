package service

import (
	"Backend-trainee-assignment/internal/model"
	"Backend-trainee-assignment/internal/model/service"
	"Backend-trainee-assignment/internal/model/transport"
	"Backend-trainee-assignment/internal/service/serverrors"
	"context"
	"fmt"
)

const userIDContextKey = "UserID"

type storage interface {
	GetUserByUsername(ctx context.Context, username string) (*service.User, error)
	GetPurchasesByUsername(ctx context.Context, username string) ([]*service.Purchase, error)
	GetCoinTransactionsByUsername(ctx context.Context, username string) ([]*service.CoinTransaction, error)
}

type Service struct {
	storage storage
}

func NewShopService(rep storage) *Service {
	return &Service{
		storage: rep,
	}
}

func (s *Service) GetInfo(ctx context.Context) (*transport.InfoResponse, error) {

	username := ctx.Value(userIDContextKey).(string)

	user, err := s.storage.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", serverrors.ErrStorage, err)
	}

	purchases, err := s.storage.GetPurchasesByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", serverrors.ErrStorage, err)
	}

	coinTransactions, err := s.storage.GetCoinTransactionsByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", serverrors.ErrStorage, err)
	}

	info := model.ParseToTransportInfo(user, purchases, coinTransactions)

	return info, nil
}

func (s *Service) SendCoins(ctx context.Context, target string) error {
	return nil
}

func (s *Service) BuyItem(ctx context.Context, id string) error {
	return nil
}

package service

import (
	"Backend-trainee-assignment/internal/model/service"
	"context"
)

const userIDContextKey = "UserID"

type storage interface {
	GetUserByUserId(ctx context.Context, userID int64) (*service.User, error)
	GetPurchasesByUserID(ctx context.Context, userID int64) ([]*service.Purchase, error)
	GetCoinTransactionsByUserID(ctx context.Context, userID int64) ([]*service.CoinTransaction, error)
}

type Service struct {
	storage storage
}

func NewShopService(rep storage) *Service {
	return &Service{
		storage: rep,
	}
}

func (s *Service) GetInfo(ctx context.Context) (string, error) {

	userID := ctx.Value(userIDContextKey).(int64)

	user, err := s.storage.GetUserByUserId(ctx, userID)
	if err != nil {
		return "", err
	}

	purchases, err := s.storage.GetPurchasesByUserID(ctx, userID)
	if err != nil {
		return "", err
	}

	coinTransactions, err := s.storage.GetCoinTransactionsByUserID(ctx, userID)
	if err != nil {
		return "", err
	}

}

func (s *Service) SendCoins(ctx context.Context, target string) error {
	return nil
}

func (s *Service) BuyItem(ctx context.Context, id string) error {
	return nil
}

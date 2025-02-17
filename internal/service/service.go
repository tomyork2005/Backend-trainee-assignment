package service

import (
	"Backend-trainee-assignment/internal/model"
	"Backend-trainee-assignment/internal/model/service"
	"Backend-trainee-assignment/internal/model/transport"
	"Backend-trainee-assignment/internal/service/serverrors"
	"context"
	"fmt"
	"log/slog"
)

const usernameContextKey = "Username"

type storage interface {
	GetUserByUsername(ctx context.Context, username string) (*service.User, error)
	GetPurchasesByUsername(ctx context.Context, username string) ([]*service.Purchase, error)
	GetCoinTransactionsByUsername(ctx context.Context, username string) ([]*service.CoinTransaction, error)
	TransferCoinsToTarget(ctx context.Context, username string, target string, amount int) error
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

	username := ctx.Value(usernameContextKey).(string)

	user, err := s.storage.GetUserByUsername(ctx, username)
	if err != nil {
		slog.Error("GetUserByUsername failed", "error", err, "username:", username)
		return nil, fmt.Errorf("%w:%w", serverrors.ErrStorage, err)
	}

	purchases, err := s.storage.GetPurchasesByUsername(ctx, username)
	if err != nil {
		slog.Error("GetPurchasesByUsername failed", "err", err, "username:", username)
		return nil, fmt.Errorf("%w:%w", serverrors.ErrStorage, err)
	}

	coinTransactions, err := s.storage.GetCoinTransactionsByUsername(ctx, username)
	if err != nil {
		slog.Error("GetCoinTransactionsByUsername failed", "error", err, "username:", username)
		return nil, fmt.Errorf("%w:%w", serverrors.ErrStorage, err)
	}

	info := model.ParseToTransportInfo(user, purchases, coinTransactions)

	slog.Info("Successfully get info for user", "username", username)

	return info, nil
}

func (s *Service) SendCoins(ctx context.Context, target string, amount int) error {
	targetUser, err := s.storage.GetUserByUsername(ctx, target)
	if err != nil {
		slog.Error("Failed to get user by username", "user", target)
		return fmt.Errorf("%w:%w", serverrors.ErrStorage, err)
	}

	if targetUser == nil {
		slog.Error("Target user not found", "target", target)
		return fmt.Errorf("%w:%w", serverrors.ErrInvalidTarget, err)
	}

	user, err := s.storage.GetUserByUsername(ctx, ctx.Value(usernameContextKey).(string))
	if err != nil {
		slog.Error("Failed to get user by username", "user", ctx.Value(usernameContextKey).(string))
		return fmt.Errorf("%w:%w", serverrors.ErrStorage, err)
	}

	if user.Balance < int64(amount) {
		slog.Error("User balance not enough", "user", ctx.Value(usernameContextKey).(string))
		return serverrors.ErrBalanceNotEnough
	}

	err = s.storage.TransferCoinsToTarget(ctx, ctx.Value(usernameContextKey).(string), target, amount)
	if err != nil {
		slog.Error("Failed to transfer coins to target", "user", ctx.Value(usernameContextKey).(string), "err", err)
		return fmt.Errorf("%w:%w", serverrors.ErrStorage, err)
	}

	slog.Info("Successfully send coins to target", "target", target, "from user ", ctx.Value(usernameContextKey).(string))

	return nil
}

func (s *Service) BuyItem(ctx context.Context, id string) error {
	return nil
}

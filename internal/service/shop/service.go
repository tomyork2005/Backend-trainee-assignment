package shop

import (
	"Backend-trainee-assignment/internal/model"
	"Backend-trainee-assignment/internal/model/service"
	"Backend-trainee-assignment/internal/model/transport"
	serviceErr "Backend-trainee-assignment/internal/service"
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
	GetMerchPrice(ctx context.Context, merchName string) (int, bool, error)
	CreatePurchase(ctx context.Context, username string, merchName string, price int) error
}

type ShopService struct {
	storage storage
}

func NewShopService(rep storage) *ShopService {
	return &ShopService{
		storage: rep,
	}
}

func (s *ShopService) GetInfo(ctx context.Context) (*transport.InfoResponse, error) {

	username := ctx.Value(usernameContextKey).(string)

	user, err := s.storage.GetUserByUsername(ctx, username)
	if err != nil {
		slog.Error("GetUserByUsername failed", "error", err, "username:", username)
		return nil, fmt.Errorf("%w:%w", serviceErr.ErrStorage, err)
	}

	purchases, err := s.storage.GetPurchasesByUsername(ctx, username)
	if err != nil {
		slog.Error("GetPurchasesByUsername failed", "err", err, "username:", username)
		return nil, fmt.Errorf("%w:%w", serviceErr.ErrStorage, err)
	}

	coinTransactions, err := s.storage.GetCoinTransactionsByUsername(ctx, username)
	if err != nil {
		slog.Error("GetCoinTransactionsByUsername failed", "error", err, "username:", username)
		return nil, fmt.Errorf("%w:%w", serviceErr.ErrStorage, err)
	}

	info := model.ParseToTransportInfo(user, purchases, coinTransactions)

	slog.Info("Successfully get info for user", "username", username)

	return info, nil
}

func (s *ShopService) SendCoins(ctx context.Context, target string, amount int) error {
	targetUser, err := s.storage.GetUserByUsername(ctx, target)
	if err != nil {
		slog.Error("Failed to get user by username", "user", target)
		return fmt.Errorf("%w:%w", serviceErr.ErrStorage, err)
	}

	if targetUser == nil {
		slog.Error("Target user not found", "target", target)
		return fmt.Errorf("%w:%w", serviceErr.ErrInvalidTarget, err)
	}

	user, err := s.storage.GetUserByUsername(ctx, ctx.Value(usernameContextKey).(string))
	if err != nil {
		slog.Error("Failed to get user by username", "user", ctx.Value(usernameContextKey).(string))
		return fmt.Errorf("%w:%w", serviceErr.ErrStorage, err)
	}

	if user.Balance < int64(amount) {
		slog.Error("User balance not enough", "user", ctx.Value(usernameContextKey).(string))
		return serviceErr.ErrBalanceNotEnough
	}

	err = s.storage.TransferCoinsToTarget(ctx, ctx.Value(usernameContextKey).(string), target, amount)
	if err != nil {
		slog.Error("Failed to transfer coins to target", "user", ctx.Value(usernameContextKey).(string), "err", err)
		return fmt.Errorf("%w:%w", serviceErr.ErrStorage, err)
	}

	slog.Info("Successfully send coins to target", "target", target, "from user ", ctx.Value(usernameContextKey).(string))

	return nil
}

func (s *ShopService) BuyItem(ctx context.Context, merchName string) error {

	price, exist, err := s.storage.GetMerchPrice(ctx, merchName)
	if err != nil {
		slog.Error("GetMerchPrice failed", "error", err)
		return fmt.Errorf("%w:%w", serviceErr.ErrStorage, err)
	}

	if !exist {
		slog.Error("GetMerchPrice not found", "merchName", merchName)
		return serviceErr.ErrMerchDontExist
	}

	user, err := s.storage.GetUserByUsername(ctx, ctx.Value(usernameContextKey).(string))
	if err != nil {
		slog.Error("Failed to get user by username", "user", ctx.Value(usernameContextKey).(string))
		return fmt.Errorf("%w:%w", serviceErr.ErrStorage, err)
	}

	if user.Balance < int64(price) {
		slog.Error("Balance does not enough", "user", ctx.Value(usernameContextKey).(string))
		return serviceErr.ErrBalanceNotEnough
	}

	err = s.storage.CreatePurchase(ctx, ctx.Value(usernameContextKey).(string), merchName, price)
	if err != nil {
		slog.Error("Failed to create purchase", "user", ctx.Value(usernameContextKey).(string), "err", err)
		return fmt.Errorf("%w:%w", serviceErr.ErrStorage, err)
	}

	slog.Info("Successfully buy purchase", "merchName", merchName, "price", price)

	return nil
}

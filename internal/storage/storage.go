package storage

import (
	"Backend-trainee-assignment/internal/model/service"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type storage struct {
	db *pgxpool.Pool
}

func NewStorage(connString string) (*storage, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &storage{
		db: pool,
	}, nil
}

func (s *storage) GetUserByUserId(ctx context.Context, userID int64) (*service.User, error) {

	row := s.db.QueryRow(ctx, queryGetUserByUserID, userID)

	var result service.User
	err := row.Scan(&result)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (s *storage) GetPurchasesByUserID(ctx context.Context, userID int64) ([]*service.Purchase, error) {

	rows, err := s.db.Query(ctx, queryGetPurchasesByUserID, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*service.Purchase, 0)
	for rows.Next() {
		var purchases service.Purchase
		err = rows.Scan(&purchases)
		if err != nil {
			return nil, err
		}
		result = append(result, &purchases)
	}

	return result, nil
}

func (s *storage) GetCoinTransactionsByUserID(ctx context.Context, userID int64) ([]*service.CoinTransaction, error) {

	rows, err := s.db.Query(ctx, queryGetCoinTransactionsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*service.CoinTransaction, 0)
	for rows.Next() {
		var coinTransaction service.CoinTransaction
		err = rows.Scan(&coinTransaction)
		if err != nil {
			return nil, err
		}
		result = append(result, &coinTransaction)
	}

	return result, nil
}

func (s *storage) GetUserByUsername(ctx context.Context, username string) (*service.User, error) {
	return nil, nil
}

func (s *storage) CreateUser(ctx context.Context, username, password string) (*service.User, error) {
	return &service.User{
		ID:       123,
		Username: username,
		Password: password,
	}, nil
}

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

func (s *storage) GetUserByUsername(ctx context.Context, username string) (*service.User, error) {

	row := s.db.QueryRow(ctx, queryGetUserByUsername, username)

	var user service.User

	err := row.Scan(
		&user.Username,
		&user.Password,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *storage) GetPurchasesByUsername(ctx context.Context, username string) ([]*service.Purchase, error) {

	rows, err := s.db.Query(ctx, queryGetPurchasesByUsername, username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*service.Purchase, 0)
	for rows.Next() {
		var purchases service.Purchase
		err = rows.Scan(
			&purchases.ID,
			&purchases.User,
			&purchases.Merch,
			&purchases.Quantity,
			&purchases.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &purchases)
	}

	return result, nil
}

func (s *storage) GetCoinTransactionsByUsername(ctx context.Context, username string) ([]*service.CoinTransaction, error) {

	rows, err := s.db.Query(ctx, queryGetCoinTransactionsByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*service.CoinTransaction, 0)
	for rows.Next() {
		var coinTransaction service.CoinTransaction
		err = rows.Scan(
			&coinTransaction.ID,
			&coinTransaction.FromUser,
			&coinTransaction.ToUser,
			&coinTransaction.Amount,
			&coinTransaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &coinTransaction)
	}

	return result, nil
}

func (s *storage) TransferCoinsToTarget(ctx context.Context, username string, target string, amount int) error {

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	var balance int
	err = tx.QueryRow(ctx, queryUpdateFrom, amount, username).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < 0 {
		return errors.New("not enough balance")
	}

	_, err = tx.Exec(ctx, queryUpdateTo, amount, target)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, queryInsertTransaction, username, target, amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) CreateUser(ctx context.Context, username, password string) (*service.User, error) {
	var user service.User

	err := s.db.QueryRow(ctx, queryCreateUser, username, password).Scan(
		&user.Username,
		&user.Password,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

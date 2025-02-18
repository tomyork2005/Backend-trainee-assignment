package storage

import (
	"Backend-trainee-assignment/internal/model/service"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(connString string) (*Storage, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: pool,
	}, nil
}

func (s *Storage) Ping() error {
	return s.db.Ping(context.Background())
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*service.User, error) {

	row := s.db.QueryRow(ctx, queryGetUserByUsername, username)

	var user service.User

	err := row.Scan(
		&user.Username,
		&user.Password,
		&user.Balance,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *Storage) GetPurchasesByUsername(ctx context.Context, username string) ([]*service.Purchase, error) {

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
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &purchases)
	}

	return result, nil
}

func (s *Storage) GetCoinTransactionsByUsername(ctx context.Context, username string) ([]*service.CoinTransaction, error) {

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
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &coinTransaction)
	}

	return result, nil
}

func (s *Storage) TransferCoinsToTarget(ctx context.Context, username string, target string, amount int) error {

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

func (s *Storage) CreateUser(ctx context.Context, username, password string) (*service.User, error) {
	var user service.User

	err := s.db.QueryRow(ctx, queryCreateUser, username, password).Scan(
		&user.Username,
		&user.Password,
		&user.Balance,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) GetMerchPrice(ctx context.Context, merchName string) (int, bool, error) {

	var price int
	err := s.db.QueryRow(ctx, queryGetMerchPrice, merchName).Scan(&price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, err
	}

	return price, true, nil
}

func (s *Storage) CreatePurchase(ctx context.Context, username string, merchName string, price int) error {

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
	err = tx.QueryRow(ctx, queryUpdateUserCoinWithReturnBalance, price, username).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < 0 {
		return errors.New("not enough balance")
	}

	_, err = tx.Exec(ctx, queryInsertPurchase, username, merchName)
	if err != nil {
		return err
	}

	return nil
}

package model

import "time"

type User struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Balance   int64     `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Merch struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Price     int       `db:"price"`
	CreatedAt time.Time `db:"created_at"`
}

type CoinTransaction struct {
	ID         int64     `db:"id"`
	FromUserID int64     `db:"from_user_id"`
	ToEUserID  int64     `db:"to_user_id"`
	Amount     int       `db:"amount"`
	CreatedAt  time.Time `db:"created_at"`
}

type Purchase struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	MerchID   int64     `db:"merch_id"`
	Quantity  int       `db:"quantity"`
	CreatedAt time.Time `db:"created_at"`
}

package service

import "time"

type User struct {
	Username  string    `db:"username"`
	Password  string    `db:"password_hash"`
	Balance   int64     `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Merch struct {
	ID        int64     `db:"id"`
	Merchname string    `db:"merchname"`
	Price     int       `db:"price"`
	CreatedAt time.Time `db:"created_at"`
}

type CoinTransaction struct {
	ID        int64     `db:"id"`
	FromUser  string    `db:"from_username"`
	ToUser    string    `db:"to_username"`
	Amount    int       `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}

type Purchase struct {
	ID        int64     `db:"id"`
	User      string    `db:"username"`
	Merch     string    `db:"merchname"`
	Quantity  int       `db:"quantity"`
	CreatedAt time.Time `db:"created_at"`
}

package service

type User struct {
	Username string `db:"username"`
	Password string `db:"password_hash"`
	Balance  int64  `db:"balance"`
}

type Merch struct {
	ID        int64  `db:"id"`
	Merchname string `db:"merchname"`
	Price     int    `db:"price"`
}

type CoinTransaction struct {
	ID       int64  `db:"id"`
	FromUser string `db:"from_username"`
	ToUser   string `db:"to_username"`
	Amount   int    `db:"amount"`
}

type Purchase struct {
	ID       int64  `db:"id"`
	User     string `db:"username"`
	Merch    string `db:"merchname"`
	Quantity int    `db:"quantity"`
}

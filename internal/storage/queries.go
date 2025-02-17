package storage

const (
	queryGetCoinTransactionsByUsername = `
	SELECT *
	FROM coin_transactions 
	WHERE from_user_id = $1 or to_user_id = $1
	`
	queryGetPurchasesByUsername = `
	SELECT *
	FROM purchases
	where username = $1
	`

	queryCreateUser = `
	INSERT INTO users(username, password_hash) VALUES ($1, $2) RETURNING username, password_hash, coins, created_at, updated_at
`
	queryGetUserByUsername = `
	SELECT username, password_hash, coins, created_at, updated_at
	FROM users WHERE username = $1
	`
)

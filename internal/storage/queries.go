package storage

const (
	queryGetCoinTransactionsByUsername = `
	SELECT id, from_user_id, to_user_id, amount, created_at
	FROM coin_transactions 
	WHERE from_user_id = $1 or to_user_id = $1
	`
	queryGetPurchasesByUsername = `
	SELECT id, username, merchname, quantity, created_at
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

	queryUpdateFrom        = `UPDATE users SET coins = coins - $1 WHERE username = $2 RETURNING coins`
	queryUpdateTo          = `UPDATE users SET coins = coins + $1 WHERE username = $2`
	queryInsertTransaction = `INSERT INTO coin_transactions (from_user_id, to_user_id, amount)
                     VALUES ($1, $2, $3)`
)

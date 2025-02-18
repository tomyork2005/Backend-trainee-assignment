package storage

// auth
const (
	queryCreateUser = `
	INSERT INTO users(username, password_hash) VALUES ($1, $2) RETURNING username, password_hash, coins
`
)

// info
const (
	queryGetCoinTransactionsByUsername = `
	SELECT id, from_user_id, to_user_id, amount
	FROM coin_transactions 
	WHERE from_user_id = $1 or to_user_id = $1
	`
	queryGetPurchasesByUsername = `
	SELECT id, username, merchname, quantity
	FROM purchases
	where username = $1
	`
	queryGetUserByUsername = `
	SELECT username, password_hash, coins
	FROM users WHERE username = $1
	`
)

// coin transaction
const (
	queryUpdateFrom        = `UPDATE users SET coins = coins - $1 WHERE username = $2 RETURNING coins`
	queryUpdateTo          = `UPDATE users SET coins = coins + $1 WHERE username = $2`
	queryInsertTransaction = `INSERT INTO coin_transactions (from_user_id, to_user_id, amount)
                     VALUES ($1, $2, $3)`
)

// buy item
const (
	queryGetMerchPrice                   = `SELECT price FROM merch WHERE merchname = $1`
	queryUpdateUserCoinWithReturnBalance = `UPDATE users SET coins = coins - $1 WHERE username = $2 RETURNING coins`
	queryInsertPurchase                  = `INSERT INTO purchases (username, merchname) VALUES ($1, $2)`
)

package storage

const (
	queryGetCoinTransactionsByUserID = `
	SELECT *
	FROM coin_transactions 
	WHERE user_id = $1
	`
	queryGetPurchasesByUserID = `
	SELECT *
	FROM purchases
	where user_id = $1
	`
	queryGetUserByUserID = `
	SELECT *
	FROM users WHERE id = $1
	`
)

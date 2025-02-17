package transport

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []Item      `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type Item struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []CoinReceived `json:"received"`
	Sent     []CoinSent     `json:"sent"`
}

type CoinReceived struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type CoinSent struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required,gte=1"`
	Amount int    `json:"amount" validate:"required,numeric"`
}

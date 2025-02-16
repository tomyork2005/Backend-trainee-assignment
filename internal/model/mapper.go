package model

import (
	"Backend-trainee-assignment/internal/model/service"
	"Backend-trainee-assignment/internal/model/transport"
)

func ParseToTransportInfo(user *service.User, purchase []*service.Purchase, transactions []*service.CoinTransaction) *transport.InfoResponse {

	var result *transport.InfoResponse

	received := make([]transport.CoinReceived, 0)
	for _, tr := range transactions {

		if tr.ToUserID == user.ID {
			received = append(received, transport.CoinReceived{
				FromUser: tr.FromUserID,
			})
		}

		coinHistory = append(coinHistory)
	}

	result = &transport.InfoResponse{
		Coins: int(user.Balance),
	}

}

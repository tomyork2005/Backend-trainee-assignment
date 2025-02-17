package model

import (
	"Backend-trainee-assignment/internal/model/service"
	"Backend-trainee-assignment/internal/model/transport"
)

func ParseToTransportInfo(user *service.User, purchases []*service.Purchase, transactions []*service.CoinTransaction) *transport.InfoResponse {

	var result *transport.InfoResponse

	countItems := make(map[string]int, 0)
	for _, pr := range purchases {
		countItems[pr.Merch] += pr.Quantity
	}

	inventory := make([]transport.Item, 0, len(countItems))
	for key, quantity := range countItems {
		inventory = append(inventory, transport.Item{
			Type:     key,
			Quantity: quantity,
		})
	}

	// don`t know how many --> prepare len == 0
	received := make([]transport.CoinReceived, 0)
	sent := make([]transport.CoinSent, 0)

	for _, tr := range transactions {
		if tr.FromUser == user.Username {
			sent = append(sent, transport.CoinSent{
				ToUser: tr.ToUser,
				Amount: tr.Amount,
			})
		} else {
			received = append(received, transport.CoinReceived{
				FromUser: tr.FromUser,
				Amount:   tr.Amount,
			})
		}
	}

	result = &transport.InfoResponse{
		Coins:     int(user.Balance),
		Inventory: inventory,
		CoinHistory: transport.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}

	return result
}

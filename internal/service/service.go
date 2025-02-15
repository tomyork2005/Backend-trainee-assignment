package service

import "context"

type Service struct {
}

func NewShopService() *Service {
	return &Service{}
}

func (s *Service) GetInfo(ctx context.Context, id string) (string, error) {
	return "", nil
}
func (s *Service) SendCoins(ctx context.Context, target string) error {
	return nil
}

func (s *Service) BuyItem(ctx context.Context, id string) error {
	return nil
}

package main

import (
	authService "Backend-trainee-assignment/internal/auth"
	"Backend-trainee-assignment/internal/config"
	shopService "Backend-trainee-assignment/internal/domain"
	"Backend-trainee-assignment/internal/transport"
	"net/http"
)

func main() {

	cfg := config.MustLoadConfig()

	auth := authService.NewAuthService(cfg.AuthConfig.TokenTTL, cfg.AuthConfig.SingingKey)
	shop := shopService.NewShopService()

	router := transport.NewHandler(shop, auth)

	r := router.Routes()

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}

}

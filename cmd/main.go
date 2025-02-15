package main

import (
	authService "Backend-trainee-assignment/internal/auth"
	shopService "Backend-trainee-assignment/internal/service"

	cfg "Backend-trainee-assignment/internal/config"
	"Backend-trainee-assignment/internal/transport"
	"net/http"
)

func main() {

	config := cfg.MustLoadConfig()

	auth := authService.NewAuthService(config.AuthConfig.TokenTTL, config.AuthConfig.SingingKey)
	shop := shopService.NewShopService()

	router := transport.NewHandler(shop, auth)

	r := router.Routes()

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}

}

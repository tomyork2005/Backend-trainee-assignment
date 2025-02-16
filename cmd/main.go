package main

import (
	authService "Backend-trainee-assignment/internal/auth"
	shopService "Backend-trainee-assignment/internal/service"
	"Backend-trainee-assignment/internal/storage"

	cfg "Backend-trainee-assignment/internal/config"
	"Backend-trainee-assignment/internal/transport"
	"net/http"
)

func main() {

	config := cfg.MustLoadConfig()

	db := storage.NewStorage()
	auth := authService.NewAuthService(config.AuthConfig.TokenTTL, config.AuthConfig.SingingKey, db)
	shop := shopService.NewShopService()

	router := transport.NewHandler(shop, auth)

	r := router.Routes()

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}

}

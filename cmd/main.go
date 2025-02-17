package main

import (
	authService "Backend-trainee-assignment/internal/auth"
	cfg "Backend-trainee-assignment/internal/config"
	shopService "Backend-trainee-assignment/internal/service"
	"Backend-trainee-assignment/internal/storage"
	"Backend-trainee-assignment/internal/transport"
	"log/slog"
	"net/http"
	"strconv"
)

func main() {

	config := cfg.MustLoadConfig()

	db, err := storage.NewStorage(postgresConnectionString(&config.PostgresConfig))
	if err != nil {
		slog.Error("Error connecting to database")
		panic(err)
	}

	auth := authService.NewAuthService(config.AuthConfig.TokenTTL, config.AuthConfig.SingingKey, db)
	shop := shopService.NewShopService(db)

	router := transport.NewHandler(shop, auth)

	r := router.Routes()

	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}

}

func postgresConnectionString(config *cfg.PostgresConfig) string {
	port := strconv.Itoa(config.Port)
	return "postgres://" + config.User + ":" + config.Password + "@" + config.Host + ":" + port + "/" + config.DbName
}

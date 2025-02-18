package app

import (
	"Backend-trainee-assignment/internal/config"
	"Backend-trainee-assignment/internal/service/auth"
	"Backend-trainee-assignment/internal/service/shop"
	postgres "Backend-trainee-assignment/internal/storage"
	"Backend-trainee-assignment/internal/transport"
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// It is possible to add interfaces taking into account that our project will grow and there will be other implementations,
// but I think in this current project is superfluous
type App struct {
	storage *postgres.Storage
	auth    *auth.AuthService
	shop    *shop.ShopService

	config *config.Config

	server *http.Server
}

func NewApp(config *config.Config) *App {

	storageInstance, err := postgres.NewStorage(postgresConnectionString(&config.PostgresConfig))
	if err != nil {
		slog.Error("Error to create storage instance", slog.String("error", err.Error()))
		panic(err)
	}

	shopServiceInstance := shop.NewShopService(storageInstance)
	authServiceInstance := auth.NewAuthService(config.AuthConfig.TokenTTL, config.AuthConfig.SingingKey, storageInstance)

	handler := transport.NewHandler(shopServiceInstance, authServiceInstance)
	router := handler.Routes()

	srv := &http.Server{
		Addr:    config.HttpConfig.Address,
		Handler: router,
	}

	return &App{
		storage: storageInstance,
		auth:    authServiceInstance,
		shop:    shopServiceInstance,

		config: config,

		server: srv,
	}
}

func (a *App) Start() error {

	err := a.storage.Ping()
	if err != nil {
		slog.Error("Error with ping storage", slog.String("error", err.Error()))
		return err
	}

	if err = a.server.ListenAndServe(); err != nil {
		slog.Error("Http server start failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (a *App) Stop() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.server.Shutdown(ctx)
	if err != nil {
		slog.Error("Http server stop failed", slog.String("error", err.Error()))
	}

	a.storage.Close()
}

func postgresConnectionString(config *config.PostgresConfig) string {
	port := strconv.Itoa(config.Port)
	return "postgres://" + config.User + ":" + config.Password + "@" + config.Host + ":" + port + "/" + config.DbName
}

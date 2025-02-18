package main

import (
	"Backend-trainee-assignment/internal/app"
	"Backend-trainee-assignment/internal/config"

	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	slog.Info("Logger with level INFO ENABLED")

	cfg := config.MustLoadConfig()
	application := app.NewApp(cfg)

	go func() {
		err := application.Start()
		if err != nil {
			os.Exit(1)
		}
	}()
	slog.Info("Application successfully started")

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
	slog.Info("Application successfully stopped")
}

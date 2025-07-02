package main

import (
	"fmt"
	"job-pulse/internal/config"
	"job-pulse/internal/storage/postgres"
	"log/slog"
	"os"
)

const (
    envLocal = "local"
    envDev   = "dev"
    envProd  = "prod"
)

func main() {
	fmt.Println("Hello, job-pulse!")
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting job-pulse", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	fmt.Println(cfg.Database)
	postgres.ConnectToDb()
	
	
}



func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

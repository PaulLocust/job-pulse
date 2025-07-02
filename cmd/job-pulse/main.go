package main

import (
	"job-pulse/internal/config"
	"job-pulse/internal/lib/sl"
	"job-pulse/internal/storage/postgres"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()
	log := sl.SetupLogger(cfg.Env)
	log.Info("starting job-pulse", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	postgres.ConnectToDb(log)
	recover()
}

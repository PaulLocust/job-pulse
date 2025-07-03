package postgres

import (
	"fmt"
	"job-pulse/backend/internal/config"
	"job-pulse/backend/internal/lib/sl"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBCon *gorm.DB

func ConnectToDb(log *slog.Logger) {
	cfg := config.MustLoad()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port, cfg.Database.SSLMode,
	)

	var err error

	DBCon, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("failed to connect database", sl.Err(err))
		os.Exit(1)
	}
	log.Info("database connected!", slog.String("db", cfg.Database.Name))

}

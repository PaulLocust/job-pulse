package postgres

import (
  "gorm.io/gorm"
  "job-pulse/internal/config"
  "gorm.io/driver/postgres"
  "fmt"
)

var DBCon *gorm.DB

func ConnectToDb() {
	cfg := config.MustLoad()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port, cfg.Database.SSLMode,
	)

	var err error

	DBCon, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

}
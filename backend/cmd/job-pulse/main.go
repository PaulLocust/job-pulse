package main

import (
	"job-pulse/backend/internal/config"
	"job-pulse/backend/internal/controllers"
	"job-pulse/backend/internal/lib/sl"
	"job-pulse/backend/internal/storage/postgres"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()
	log := sl.SetupLogger(cfg.Env)
	log.Info("starting job-pulse", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	postgres.ConnectToDb(log)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/vacancies", controllers.GetVacancies(log))

	log.Info("server started at :8080")
	r.Run(":8080")

}

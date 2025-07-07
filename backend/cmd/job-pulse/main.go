package main

import (
	"job-pulse/backend/internal/config"
	"job-pulse/backend/internal/controllers"
	"job-pulse/backend/internal/lib/sl"
	"job-pulse/backend/internal/storage/postgres"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()
	log := sl.SetupLogger(cfg.Env)
	log.Info("starting job-pulse", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	postgres.ConnectToDb(log)

	r := gin.Default()

	// Настройка CORS (разрешаем запросы с фронтенда)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // URL вашего фронтенда
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/vacancies", controllers.GetVacancies(log))
	r.GET("/api/skills/stats", controllers.GetSkillStats(log))

	log.Info("server started at :8080")
	r.Run(":8080")

}

package main

import (
	"net/http"
	"github.com/robfig/cron/v3"
	"job-pulse/backend/internal/config"
	"job-pulse/backend/internal/controllers"
	"job-pulse/backend/internal/lib/sl"
	"job-pulse/backend/internal/storage/postgres"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func setupCron(log *slog.Logger) {
	c := cron.New()

	// Обновляем данные каждые 6 часов
	c.AddFunc("0 */6 * * *", func() {
		log.Info("Starting scheduled vacancies update")

		// Вызываем наш парсер
		_, err := http.Get("http://localhost:8080/vacancies?q=golang")
		if err != nil {
			log.Error("Scheduled update failed", sl.Err(err))
		}
	})

	c.Start()
}

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
	r.GET("/vacancies/:id", controllers.GetVacancyDetails(log))
	r.GET("/vacancies/tech/:id", controllers.GetVacancyTechDetails(log))

	setupCron(log)

	log.Info("server started at :8080")
	r.Run(":8080")
	
}

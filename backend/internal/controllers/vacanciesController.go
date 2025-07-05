package controllers

import (
	"job-pulse/backend/internal/hhapi"
	"job-pulse/backend/internal/lib/sl"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetVacancies(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			query = "golang"
		}

		// Получаем и обрабатываем вакансии
		vacancies, err := hhapi.FetchAndProcessVacancies(query, log)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch and process vacancies"})
			return
		}

		// Сохраняем вакансии
		savedCount := 0
		for _, vacancy := range vacancies {
			if err := hhapi.SaveVacancy(vacancy); err != nil {
				log.Error("failed to save vacancy", 
					slog.String("id", vacancy.ID), 
					sl.Err(err))
				continue
			}
			savedCount++
		}

		c.JSON(http.StatusOK, gin.H{
			"items":       vacancies,
			"saved_count": savedCount,
		})
	}
}
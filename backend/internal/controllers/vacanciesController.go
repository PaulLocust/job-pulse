package controllers

import (
	"job-pulse/backend/internal/hhapi"
	"job-pulse/backend/internal/lib/sl"
	"job-pulse/backend/internal/lib/statistics"
	"log/slog"
	"net/http"
	"strconv"

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

func GetSkillStats(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "20") // Значение по умолчанию: 20
		limit, err := strconv.Atoi(limitStr)
		
		if limit > 100 {
    		limit = 100
		}

		if err != nil {
    		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a number"})
    		return
		}

    	stats, err := statistics.GetSkillsFrequency(limit)
    	if err != nil {
    		log.Error("Failed to get skills stats", "error", err)
    		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
    		return
		}

    	c.JSON(http.StatusOK, stats)
	}
}
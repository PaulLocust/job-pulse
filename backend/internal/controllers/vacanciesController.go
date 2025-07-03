package controllers

import (
	"job-pulse/backend/internal/hhapi"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetVacancies(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			query = "golang"
		}

		vacanciesResp, err := hhapi.FetchVacancies(query, log)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch vacancies"})
			return
		}

		// Учёт двоякого написания языка в вакансии, пример: (Go/Golang)
		langVariants := []string{}
		switch queryLower := strings.ToLower(query); queryLower {
		case "golang", "go":
			langVariants = []string{"go", "golang"}
		default:
			langVariants = []string{query}
		}

		excludeLangs := []string{"php", "python", "java", "c++", "javascript", "ruby", "golang", "go"}

		filteredItems := hhapi.FilterVacanciesByLanguages(vacanciesResp.Items, langVariants, excludeLangs)

		c.JSON(http.StatusOK, gin.H{
			"items": filteredItems,
		})
	}
}

func GetVacancyDetails(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vacancy id required"})
			return
		}

		details, err := hhapi.FetchVacancyDetails(id, log)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch vacancy details"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":          details.ID,
			"name":        details.Name,
			"description": details.Description,
			"key_skills":  details.KeySkills,
			"salary":      details.Salary,
		})
	}
}

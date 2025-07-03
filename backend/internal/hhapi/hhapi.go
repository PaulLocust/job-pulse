package hhapi

import (
	"encoding/json"
	"job-pulse/backend/internal/lib/sl"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Функции для запросов на hh.ru

func fetchVacancies(query string, log *slog.Logger) (*VacanciesResponse, error) {
	baseURL := "https://api.hh.ru/vacancies"
	params := url.Values{}
	params.Add("text", query)
	params.Add("per_page", "20")
	

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		log.Error("failed to request hh.ru API", sl.Err(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("hh.ru API returned non-200 status", sl.Status(resp))
		return nil, err
	}

	var vacancies VacanciesResponse
	if err := json.NewDecoder(resp.Body).Decode(&vacancies); err != nil {
		log.Error("failed to decode hh.ru API response", sl.Err(err))
		return nil, err
	}

	return &vacancies, nil
}

func GetVacancies(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Новый эндпоинт для проверки получения вакансий
		query := c.Query("q")
		if query == "" {
			query = "Golang" // значение по умолчанию
		}

		vacancies, err := fetchVacancies(query, log)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch vacancies"})
			return
		}

		c.JSON(http.StatusOK, vacancies)
	}
}

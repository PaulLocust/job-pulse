package hhapi

import (
	"encoding/json"
	"fmt"
	"job-pulse/backend/internal/lib/sl"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
)

func FetchVacancies(query string, log *slog.Logger) (*VacanciesResponse, error) {
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

func FetchVacancyDetails(id string, log *slog.Logger) (*VacancyDetails, error) {
	url := "https://api.hh.ru/vacancies/" + id
	resp, err := http.Get(url)
	if err != nil {
		log.Error("failed to request vacancy details", sl.Err(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("vacancy details API returned non-200 status", sl.Status(resp))
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var details VacancyDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		log.Error("failed to decode vacancy details response", sl.Err(err))
		return nil, err
	}

	return &details, nil
}

func FilterVacanciesByLanguages(vacancies []Vacancy, langVariants []string, excludeLangs []string) []Vacancy {
	// Компилируем регулярные выражения для вариантов искомого языка
	var langPatterns []*regexp.Regexp
	for _, lang := range langVariants {
		langPatterns = append(langPatterns, regexp.MustCompile(`(?i)\b`+regexp.QuoteMeta(lang)+`\b`))
	}

	// Компилируем паттерны для исключаемых языков
	var excludePatterns []*regexp.Regexp
	for _, l := range excludeLangs {
		// пропускаем все варианты искомого языка, чтобы не исключать их
		skip := false
		for _, v := range langVariants {
			if l == v {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		excludePatterns = append(excludePatterns, regexp.MustCompile(`(?i)\b`+regexp.QuoteMeta(l)+`\b`))
	}

	var filtered []Vacancy
	for _, v := range vacancies {
		title := v.Name

		// Проверяем, есть ли хотя бы один вариант языка в названии
		matchedLang := false
		for _, p := range langPatterns {
			if p.MatchString(title) {
				matchedLang = true
				break
			}
		}
		if !matchedLang {
			continue
		}
		
		// Проверяем, нет ли в названии исключаемых языков
		excluded := false
		for _, p := range excludePatterns {
			if p.MatchString(title) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		filtered = append(filtered, v)
	}
	return filtered
}

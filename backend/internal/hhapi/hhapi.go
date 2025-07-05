package hhapi

import (
	"encoding/json"
	"fmt"
	"job-pulse/backend/internal/lib/dataset"
	"job-pulse/backend/internal/lib/sl"
	"job-pulse/backend/internal/storage/postgres"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	strip "github.com/grokify/html-strip-tags-go"
)

const baseURL = "https://api.hh.ru/vacancies"
const datasetPath = "D:/my-repo/Go-projects/job-pulse/backend/internal/lib/dataset/tech_dataset.json"

// SaveVacancy сохраняет вакансию в базу данных с проверкой на nil и поврежденную кодировку
func SaveVacancy(vacancy VacancyProcessed) error {
	// Проверяем существование вакансии с помощью Find
	var existing postgres.Vacancy
	result := postgres.DBCon.
		Where("hh_id = ?", vacancy.ID).
		Limit(1).
		Find(&existing)

	if result.Error != nil {
		slog.Error("Ошибка при проверке вакансии",
			slog.String("id", vacancy.ID),
			slog.String("error", result.Error.Error()))
		return result.Error
	}

	// Если вакансия уже существует - пропускаем сохранение
	if result.RowsAffected > 0 {
		return nil
	}

	// Создаем базовую запись вакансии
	v := postgres.Vacancy{
		HHID: vacancy.ID,
		Name: vacancy.Name,
	}

	// Обрабатываем salary, если он есть
	if vacancy.Salary != nil {
		v.SalaryFrom = vacancy.Salary.From
		v.SalaryTo = vacancy.Salary.To
		v.Currency = vacancy.Salary.Currency
	}

	// Обрабатываем опыт работы, если он есть
	if vacancy.Experience != nil {
		v.ExpID = vacancy.Experience.ID
	}

	// Сохраняем вакансию
	if err := postgres.DBCon.Create(&v).Error; err != nil {
		slog.Error("Ошибка сохранения вакансии",
			slog.String("id", vacancy.ID),
			slog.String("error", err.Error()))
		return err
	}

	// Сохраняем навыки
	for _, skillName := range vacancy.Skills {
		// Нормализуем название навыка
		skillName = strings.TrimSpace(skillName)
		if skillName == "" {
			continue
		}

		// Пропускаем навык, если его название содержит невалидные UTF-8 символы
		if !utf8.ValidString(skillName) {
			slog.Warn("Пропуск навыка с поврежденной кодировкой",
				slog.String("vacancy_id", vacancy.ID),
				slog.String("skill", skillName))
			continue
		}

		// Проверяем существование навыка с помощью Find
		var skill postgres.Skill
		skillResult := postgres.DBCon.
			Where("LOWER(name) = LOWER(?)", skillName).
			Limit(1).
			Find(&skill)

		if skillResult.Error != nil {
			slog.Error("Ошибка при проверке навыка",
				slog.String("skill", skillName),
				slog.String("error", skillResult.Error.Error()))
			continue
		}

		// Если навык не найден - создаем новый
		if skillResult.RowsAffected == 0 {
			skill = postgres.Skill{Name: skillName}
			if err := postgres.DBCon.Create(&skill).Error; err != nil {
				slog.Error("Ошибка создания навыка",
					slog.String("skill", skillName),
					slog.String("error", err.Error()))
				continue
			}
		}

		// Связываем навык с вакансией
		if err := postgres.DBCon.Model(&v).Association("Skills").Append(&skill); err != nil {
			slog.Error("Ошибка связывания навыка с вакансией",
				slog.String("vacancy_id", vacancy.ID),
				slog.String("skill", skillName),
				slog.String("error", err.Error()))
		}
	}

	return nil
}

func FetchAndProcessVacancies(query string, log *slog.Logger) ([]VacancyProcessed, error) {
	// 1. Получаем базовые данные вакансий
	vacancies, err := fetchVacancies(query, log)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vacancies: %w", err)
	}

	// 2. Фильтруем по языкам
	filtered := filterVacancies(vacancies, query, log)

	// 3. Загружаем dataset технологий
	techSet, err := dataset.LoadTechDataset(datasetPath)
	if err != nil {
		log.Error("failed to load tech dataset", sl.Err(err))
		techSet = make(map[string]bool)
	}

	// 4. Обрабатываем каждую вакансию
	var result []VacancyProcessed
	for _, v := range filtered {
		processedVacancy, err := processVacancy(v.ID, techSet, log)
		if err != nil {
			log.Error("failed to process vacancy", slog.String("id", v.ID), sl.Err(err))
			continue
		}
		result = append(result, *processedVacancy)
	}

	return result, nil
}

func fetchVacancies(query string, log *slog.Logger) ([]BasicVacancy, error) {
	params := url.Values{}
	params.Add("text", query)
	params.Add("per_page", "100")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var data VacanciesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Items, nil
}

func filterVacancies(vacancies []BasicVacancy, query string, log *slog.Logger) []BasicVacancy {
	langVariants := getLanguageVariants(query)
	excludeLangs := []string{"php", "python", "java", "c++", "javascript", "ruby"}

	langPatterns := compilePatterns(langVariants)
	excludePatterns := compileExcludePatterns(langVariants, excludeLangs)

	var filtered []BasicVacancy
	for _, v := range vacancies {
		if !matchesAny(v.Name, langPatterns) {
			continue
		}
		if matchesAny(v.Name, excludePatterns) {
			continue
		}
		filtered = append(filtered, v)
	}

	return filtered
}

func processVacancy(id string, techSet map[string]bool, log *slog.Logger) (*VacancyProcessed, error) {
	details, err := fetchVacancyDetails(id, log)
	if err != nil {
		return nil, err
	}

	return &VacancyProcessed{
		ID:          details.ID,
		Name:        details.Name,
		Description: details.Description,
		Salary:      details.Salary,
		Experience:  details.Experience,
		Skills:      extractSkills(details.KeySkills, details.Description, techSet),
	}, nil
}

// Вспомогательные функции
func getLanguageVariants(query string) []string {
	query = strings.ToLower(query)
	if query == "golang" || query == "go" {
		return []string{"go", "golang"}
	}
	return []string{query}
}

func compilePatterns(langs []string) []*regexp.Regexp {
	var patterns []*regexp.Regexp
	for _, lang := range langs {
		patterns = append(patterns, regexp.MustCompile(`(?i)\b`+regexp.QuoteMeta(lang)+`\b`))
	}
	return patterns
}

func compileExcludePatterns(langVariants, excludeLangs []string) []*regexp.Regexp {
	var patterns []*regexp.Regexp
	for _, lang := range excludeLangs {
		skip := false
		for _, variant := range langVariants {
			if lang == variant {
				skip = true
				break
			}
		}
		if !skip {
			patterns = append(patterns, regexp.MustCompile(`(?i)\b`+regexp.QuoteMeta(lang)+`\b`))
		}
	}
	return patterns
}

func matchesAny(text string, patterns []*regexp.Regexp) bool {
	for _, p := range patterns {
		if p.MatchString(text) {
			return true
		}
	}
	return false
}

func fetchVacancyDetails(id string, log *slog.Logger) (*VacancyRaw, error) {
	url := fmt.Sprintf("%s/%s", baseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vacancy API returned status %d", resp.StatusCode)
	}

	var details VacancyRaw
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}

func extractSkills(keySkills []Skill, description string, techSet map[string]bool) []string {
	skillsMap := make(map[string]bool)

	// 1. Обрабатываем ключевые навыки
	for _, skill := range keySkills {
		normalized := normalizeSkillName(skill.Name)
		if normalized != "" {
			skillsMap[normalized] = true
		}
	}

	// 2. Извлекаем навыки из описания
	cleanDesc := strip.StripTags(description)
	for tech := range techSet {
		if strings.Contains(strings.ToLower(cleanDesc), strings.ToLower(tech)) {
			normalized := normalizeSkillName(tech)
			if normalized != "" {
				skillsMap[normalized] = true
			}
		}
	}

	// 3. Конвертируем в слайс и сортируем
	result := make([]string, 0, len(skillsMap))
	for skill := range skillsMap {
		result = append(result, skill)
	}
	sort.Strings(result)

	return result
}

func normalizeSkillName(skill string) string {
	// Приводим к нижнему регистру и обрезаем пробелы
	skill = strings.TrimSpace(strings.ToLower(skill))

	// Удаляем лишние слова
	skill = regexp.MustCompile(`(?i)\b(знание|опыт|работа с|умение)\b`).ReplaceAllString(skill, "")
	skill = strings.TrimSpace(skill)

	// Заменяем алиасы
	if alias, ok := skillAliases[skill]; ok {
		return alias
	}

	// Стандартизируем написание
	switch {
	case strings.Contains(skill, "go"):
		return "Go"
	case strings.Contains(skill, "postgres"):
		return "PostgreSQL"
	case strings.Contains(skill, "kubernetes"):
		return "Kubernetes"
	case strings.Contains(skill, "ci/cd"), strings.Contains(skill, "ci cd"):
		return "CI/CD"
	case strings.Contains(skill, "rest"):
		return "REST API"
	case strings.Contains(skill, "grpc"):
		return "gRPC"
	case strings.Contains(skill, "unit test"), strings.Contains(skill, "unit-тест"):
		return "Unit Testing"
	case strings.Contains(skill, "микросервис"), strings.Contains(skill, "микросервисы"):
		return "Microservices"
	}

	// Первая буква заглавная для остальных
	if len(skill) > 0 {
		return strings.ToUpper(skill[:1]) + skill[1:]
	}

	return ""
}

var skillAliases = map[string]string{
	"golang":       "Go",
	"go lang":      "Go",
	"postgresql":   "PostgreSQL",
	"postgres":     "PostgreSQL",
	"k8s":          "Kubernetes",
	"ci/cd":        "CI/CD",
	"rest api":     "REST API",
	"rest":         "REST API",
	"unit-тест":    "Unit Testing",
	"unit testing": "Unit Testing",
	"nosql":        "NoSQL",
	"tcp/ip":       "TCP/IP",
	"gitlab ci":    "GitLab CI",
	"gitlab":       "GitLab",
	"алгоритмы":    "Algorithms",
	"микросервисная архитектура": "Microservices",
	"микросервисы":               "Microservices",
}

package hhapi

import (
	"job-pulse/backend/internal/lib/dataset"
	"job-pulse/backend/internal/lib/sl"
	"job-pulse/backend/internal/storage/postgres"
	"log/slog"
	"regexp"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
)

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

var techAliases = map[string]string{
	"Postgres": "PostgreSQL",
	"k8s":      "Kubernetes",
	"gRPC":     "GRPC",
	"Golang":   "Go",
}

func NormalizeTechName(tech string) string {
	tech = strings.ToLower(tech)
	for alias, normalized := range techAliases {
		if strings.Contains(tech, strings.ToLower(alias)) {
			return normalized
		}
	}
	return tech
}

// GetVacancyTechDetails - основная функция для получения данных с технологиями
func GetVacancyTechDetails(id string, log *slog.Logger, datasetPath string) (*VacancyTech, error) {
	details, err := FetchVacancyDetails(id, log)
	if err != nil {
		return nil, err
	}

	techSet, err := dataset.LoadTechDataset(datasetPath)
	if err != nil {
		log.Error("failed to load tech dataset", sl.Err(err))
		techSet = make(map[string]bool)
	}

	skills := extractAllTechnologies(details, techSet)

	return &VacancyTech{
		ID:         details.ID,
		Name:       details.Name,
		Salary:     details.Salary,
		Skills:     skills,
		Experience: details.Experience, // Добавляем опыт работы
	}, nil
}

// extractAllTechnologies объединяет технологии из key_skills и description
func extractAllTechnologies(details *VacancyDetails, techSet map[string]bool) []string {
	// 1. Технологии из key_skills
	skillsMap := make(map[string]bool)
	for _, skill := range details.KeySkills {
		normalized := NormalizeTechName(skill.Name)
		skillsMap[normalized] = true
	}

	// 2. Технологии из description
	cleanDesc := strip.StripTags(details.Description)
	for tech := range techSet {
		if strings.Contains(cleanDesc, tech) {
			normalized := NormalizeTechName(tech)
			skillsMap[normalized] = true
		}
	}

	// Преобразуем map в slice
	result := make([]string, 0, len(skillsMap))
	for tech := range skillsMap {
		result = append(result, tech)
	}

	return result
}

func SaveVacancy(vacancy Vacancy, log *slog.Logger) error {
	details, err := FetchVacancyDetails(vacancy.ID, log)
	if err != nil {
		return err
	}

	// Конвертируем в модель БД
	v := postgres.Vacancy{
		HHID:       details.ID,
		Name:       details.Name,
		SalaryFrom: details.Salary.From,
		SalaryTo:   details.Salary.To,
		Currency:   details.Salary.Currency,
		ExpID:      details.Experience.ID,
	}

	// Обрабатываем навыки
	for _, skill := range details.KeySkills {
		s := postgres.Skill{Name: skill.Name}
		if err := postgres.DBCon.FirstOrCreate(&s, "name = ?", s.Name).Error; err != nil {
			log.Error("failed to process skill", slog.String("skill", skill.Name), sl.Err(err))
			continue
		}
		v.Skills = append(v.Skills, s)
	}

	// Сохраняем вакансию
	return postgres.DBCon.Create(&v).Error
}

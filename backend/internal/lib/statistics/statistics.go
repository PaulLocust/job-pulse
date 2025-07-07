package statistics

import (
	"fmt"
	"job-pulse/backend/internal/storage/postgres"
	"log"
)

// SkillStats представляет статистику по навыку
type SkillStats struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GetSkillsFrequency возвращает статистику по частоте встречаемости навыков
func GetSkillsFrequency(limit int) ([]SkillStats, error) {
	var stats []SkillStats

	// Запрос для подсчета частоты навыков
	err := postgres.DBCon.
		Table("skills").
		Select("skills.name, COUNT(vacancy_skills.skill_id) as count").
		Joins("LEFT JOIN vacancy_skills ON skills.id = vacancy_skills.skill_id").
		Group("skills.name").
		Order("count DESC").
		Limit(limit).
		Find(&stats).Error

	if err != nil {
		return nil, err
	}

	return stats, nil
}

func PrintSkillsStats() {
	stats, err := GetSkillsFrequency(20) // Топ-20 навыков
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
	
	fmt.Println("Топ навыков в вакансиях:")
	fmt.Println("------------------------")
	for i, stat := range stats {
		fmt.Printf("%d. %s: %d упоминаний\n", i+1, stat.Name, stat.Count)
	}
}

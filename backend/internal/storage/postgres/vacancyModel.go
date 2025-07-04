package postgres

import "gorm.io/gorm"

type Vacancy struct {
	gorm.Model
	HHID       string   `gorm:"uniqueIndex;size:20"` // ID вакансии из HH.ru
	Name       string   `gorm:"not null"`
	SalaryFrom *int     `gorm:"default:null"`
	SalaryTo   *int     `gorm:"default:null"`
	Currency   string   `gorm:"size:3;default:'RUR'"`
	ExpID      string   `gorm:"size:20"`  // between1And3, moreThan3 и т.д.
	Skills     []Skill  `gorm:"many2many:vacancy_skills;"`
}

type Skill struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"unique;size:100;collate:case_insensitive"` 
}
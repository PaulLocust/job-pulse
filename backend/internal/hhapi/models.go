package hhapi

// Структуры вакансий и ответов

type Salary struct {
	From     *int   `json:"from"`
	To       *int   `json:"to"`
	Currency string `json:"currency"`
	Gross    bool   `json:"gross"`
}

type Skill struct {
	Name string `json:"name"`
}

type VacancyDetails struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	KeySkills   []Skill `json:"key_skills"` // Это массив объектов Skill, а не строка
	Salary      *Salary `json:"salary"`
	// другие поля при необходимости
}

type Vacancy struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Salary *Salary `json:"salary"`
	// другие поля вакансии
}

type VacanciesResponse struct {
	Items []Vacancy `json:"items"`
	// другие поля ответа
}

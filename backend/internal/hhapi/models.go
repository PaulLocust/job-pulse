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


type Experience struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type VacancyDetails struct {
    ID          string      `json:"id"`
    Name        string      `json:"name"`
    Description string      `json:"description"`
    KeySkills   []Skill     `json:"key_skills"`
    Salary      *Salary     `json:"salary"`
    Experience  *Experience `json:"experience"` // Добавляем опыт работы
}

type VacancyTech struct {
    ID         string     `json:"id"`
    Name       string     `json:"name"`
    Salary     *Salary    `json:"salary"`
    Skills     []string   `json:"skills"`
    Experience *Experience `json:"experience"` // Добавляем в выходную структуру
}

type Vacancy struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Salary *Salary `json:"salary"`
}

type VacanciesResponse struct {
	Items []Vacancy `json:"items"`
}

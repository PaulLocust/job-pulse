package hhapi

type Salary struct {
	From     *int   `json:"from"`
	To       *int   `json:"to"`
	Currency string `json:"currency"`
	Gross    bool   `json:"gross"`
}

type Experience struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Skill struct {
	Name string `json:"name"`
}

// BasicVacancy содержит минимальные данные о вакансии для первоначальной фильтрации
type BasicVacancy struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// VacancyRaw содержит необходимые необработанные данные вакансии из API HH
type VacancyRaw struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Salary      *Salary     `json:"salary"`
	Experience  *Experience `json:"experience"`
	KeySkills   []Skill     `json:"key_skills"`
}

// VacancyProcessed содержит обработанные данные вакансии для конечного хранения в бд
type VacancyProcessed struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Salary      *Salary     `json:"salary"`
	Skills      []string    `json:"skills"`
	Experience  *Experience `json:"experience"`
}

// VacanciesResponse представляет ответ API HH при запросе списка вакансий
type VacanciesResponse struct {
	Items []BasicVacancy `json:"items"`
}
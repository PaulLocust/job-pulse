package hhapi

// Структуры вакансий и ответов

type Salary struct {
    From     *int    `json:"from"`
    To       *int    `json:"to"`
    Currency string  `json:"currency"`
    Gross    bool    `json:"gross"`
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

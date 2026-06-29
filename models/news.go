package models

// NewsFullDetailed — полная модель новости
type NewsFullDetailed struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
}

// NewsShortDetailed — краткая модель новости для списков
type NewsShortDetailed struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

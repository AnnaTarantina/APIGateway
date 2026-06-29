package handlers

import (
	"net/http"

	"github.com/AnnaTarantina/APIGateway/models"
)

// GetNewsList возвращает список новостей (заглушка)
func GetNewsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	newsList := []models.NewsShortDetailed{
		{ID: "1", Title: "Заголовок первой новости", Summary: "Краткое содержание первой новости"},
		{ID: "2", Title: "Заголовок второй новости", Summary: "Краткое содержание второй новости"},
	}
	writeJSON(w, http.StatusOK, newsList)
}

// FilterNews возвращает отфильтрованные новости (заглушка)
func FilterNews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filteredNews := []models.NewsShortDetailed{
		{ID: "1", Title: "Отфильтрованная новость", Summary: "Отфильтрованное содержание"},
	}
	writeJSON(w, http.StatusOK, filteredNews)
}

// GetDetailedNews возвращает детальную новость (заглушка)
func GetDetailedNews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	detailedNews := models.NewsFullDetailed{
		ID:        "1",
		Title:     "Детальная новость",
		Content:   "Полный текст новости",
		Author:    "Автор",
		CreatedAt: "2024-01-01T00:00:00Z",
	}
	writeJSON(w, http.StatusOK, detailedNews)
}

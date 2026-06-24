package handlers

import (
	"encoding/json"
	"net/http"

	"../models"
)

func GetNewsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	newsList := []models.NewsShortDetailed{
		{ID: "1", Title: "Заголовок первой новости", Summary: "Краткое содержание первой новости"},
		{ID: "2", Title: "Заголовок второй новости", Summary: "Краткое содержание второй новости"},
	}

	json.NewEncoder(w).Encode(newsList)
}

func FilterNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	filteredNews := []models.NewsShortDetailed{
		{ID: "1", Title: "Отфильтрованная новость", Summary: "Отфильтрованное содержание"},
	}

	json.NewEncoder(w).Encode(filteredNews)
}

func GetDetailedNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	detailedNews := models.NewsFullDetailed{
		ID:        "1",
		Title:     "Детальная новость",
		Content:   "Полный текст новости",
		Author:    "Автор",
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	json.NewEncoder(w).Encode(detailedNews)
}

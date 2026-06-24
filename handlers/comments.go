package handlers

import (
	"encoding/json"
	"net/http"

	"../models"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	comment := models.Comment{
		ID:         "1",
		NewsID:     "1",
		Text:       "Текст комментария",
		Author:     "Пользователь",
		CreatedAt:  "2024-01-01T00:00:00Z",
		IsApproved: true,
	}

	json.NewEncoder(w).Encode(comment)
}

func GetCommentsByNewsID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	comments := []models.Comment{
		{
			ID:         "1",
			NewsID:     "1",
			Text:       "Первый комментарий",
			Author:     "Автор 1",
			CreatedAt:  "2024-01-01T00:01:00Z",
			IsApproved: true,
		},
		{
			ID:         "2",
			NewsID:     "1",
			Text:       "Второй комментарий",
			ParentID:   "1",
			Author:     "Автор 2",
			CreatedAt:  "2024-01-01T00:02:00Z",
			IsApproved: true,
		},
	}

	json.NewEncoder(w).Encode(comments)
}

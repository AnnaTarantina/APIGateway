package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Адрес микросервиса комментариев
const commentServiceURL = "http://localhost:3000"

// CommentsHandler — универсальный обработчик /comments
// маршрутизирует запросы по методу (POST — создание, GET — список)
func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddComment(w, r)
	case http.MethodGet:
		GetCommentsByNewsID(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// AddComment проксирует POST-запрос в CommentService
func AddComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Формируем запрос к CommentService
	targetURL := commentServiceURL + "/comment"
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling CommentService: %v", err)
		http.Error(w, "Comment service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Копируем статус и тело ответа
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response: %v", err)
	}
}

// GetCommentsByNewsID проксирует GET-запрос в CommentService
func GetCommentsByNewsID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	newsID := r.URL.Query().Get("news_id")
	if newsID == "" {
		http.Error(w, "news_id is required", http.StatusBadRequest)
		return
	}

	// Формируем запрос к CommentService
	targetURL := commentServiceURL + "/comments?news_id=" + newsID
	resp, err := http.Get(targetURL)
	if err != nil {
		log.Printf("Error calling CommentService: %v", err)
		http.Error(w, "Comment service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Копируем ответ
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response: %v", err)
	}
}

// writeJSON — вспомогательная функция для записи JSON
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AnnaTarantina/APIGateway/middleware"
)

const commentServiceURL = "http://localhost:8081"

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

func AddComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID, _ := ctx.Value(middleware.RequestIDKey).(string)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to read body"})
		return
	}
	defer r.Body.Close()

	// Парсим для проверки цензуры
	var commentReq struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(body, &commentReq); err == nil {
		if !IsAllowed(commentReq.Text) {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "comment contains forbidden words",
			})
			return
		}
	}

	// Проксируем в CommentService
	targetURL := fmt.Sprintf("%s/comment?request_id=%s", commentServiceURL, reqID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", reqID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "comment service unavailable"})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func GetCommentsByNewsID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID, _ := ctx.Value(middleware.RequestIDKey).(string)
	newsID := r.URL.Query().Get("news_id")

	if newsID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "news_id is required"})
		return
	}

	targetURL := fmt.Sprintf("%s/comments?news_id=%s&request_id=%s", commentServiceURL, newsID, reqID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	req.Header.Set("X-Request-ID", reqID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "comment service unavailable"})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
